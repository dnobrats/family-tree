// admin.js — shared behaviour for the admin pages: delete actions (via fetch +
// CSRF header), print buttons, and the living/dead field toggle on the person
// form. Keeping this external satisfies the Content-Security-Policy (no inline
// scripts).
(function () {
  "use strict";

  function csrfToken() {
    var m = document.cookie.match(/(?:^|;\s*)csrf_token=([^;]+)/);
    return m ? decodeURIComponent(m[1]) : "";
  }

  function postDelete(url, confirmText) {
    if (!window.confirm(confirmText)) return;
    fetch(url, { method: "POST", headers: { "X-CSRF-Token": csrfToken() } })
      .then(function (res) {
        if (!res.ok) {
          return res.text().then(function (t) { throw new Error(t || "Không thể xóa"); });
        }
        window.location.reload();
      })
      .catch(function (err) { window.alert(err.message || "Không thể xóa."); });
  }

  document.addEventListener("click", function (e) {
    var personBtn = e.target.closest("[data-delete-person]");
    if (personBtn) {
      var name = personBtn.getAttribute("data-name") || "";
      postDelete("/admin/persons/" + personBtn.getAttribute("data-delete-person") + "/delete",
        'Xóa "' + name + '"?\nThao tác này không thể hoàn tác.');
      return;
    }
    var spouseBtn = e.target.closest("[data-delete-spouse]");
    if (spouseBtn) {
      postDelete("/admin/spouses/" + spouseBtn.getAttribute("data-delete-spouse") + "/delete",
        "Xóa quan hệ vợ/chồng này?");
    }
  });

  document.querySelectorAll("[data-print]").forEach(function (btn) {
    btn.addEventListener("click", function () { window.print(); });
  });

  // Confirmation for plain form submits (CSP forbids inline onsubmit handlers).
  document.addEventListener("submit", function (e) {
    var form = e.target.closest("[data-confirm]");
    if (form && !window.confirm(form.getAttribute("data-confirm"))) {
      e.preventDefault();
    }
  });

  var alive = document.getElementById("is_alive");
  var deathFields = document.getElementById("death_fields");
  if (alive && deathFields) {
    var sync = function () { deathFields.hidden = alive.checked; };
    alive.addEventListener("change", sync);
    sync();
  }
})();
