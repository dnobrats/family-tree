// tree.js — renders the nested family tree returned by /api/tree and
// /api/clans/{id}/tree, and shows a person popup on click.
//
// The tree structure comes fully built from the server, so this file only
// renders it. No inline script is allowed (CSP), so all handlers are attached
// here.
(function () {
  "use strict";

  var state = { query: "", scale: 1 };
  var treeEl = null;

  function el(tag, cls, text) {
    var e = document.createElement(tag);
    if (cls) e.className = cls;
    if (text != null) e.textContent = text;
    return e;
  }

  function personCard(person, isSpouse) {
    var cls = "node";
    if (person.gender === 2 || person.gender === 0) cls += " female";
    if (!person.is_alive) cls += " dead";
    if (isSpouse) cls += " spouse";

    var btn = el("button", cls, person.full_name);
    btn.type = "button";
    btn.addEventListener("click", function () { showDetail(person.id); });
    return btn;
  }

  function renderNode(node) {
    var li = el("li");
    var couple = el("div", "couple");
    couple.appendChild(personCard(node.person, false));
    if (node.spouse) couple.appendChild(personCard(node.spouse, true));

    var hasChildren = node.children && node.children.length > 0;
    if (hasChildren) {
      var toggle = el("button", "toggle", "\u2212");
      toggle.type = "button";
      toggle.setAttribute("aria-label", "Thu gọn nhánh");
      toggle.addEventListener("click", function (ev) {
        ev.stopPropagation();
        var collapsed = li.classList.toggle("collapsed");
        toggle.textContent = collapsed ? "+" : "\u2212";
      });
      couple.appendChild(toggle);
    }
    li.appendChild(couple);

    if (hasChildren) {
      var ul = el("ul", "children");
      node.children.forEach(function (child) { ul.appendChild(renderNode(child)); });
      li.appendChild(ul);
    }
    return li;
  }

  function renderTree(root) {
    treeEl.innerHTML = "";
    if (!root) {
      treeEl.appendChild(el("p", "muted", "Chưa có dữ liệu."));
      return;
    }
    var ul = el("ul", "root");
    ul.appendChild(renderNode(root));
    treeEl.appendChild(ul);
    applySearch();
  }

  function applySearch() {
    var q = state.query.trim().toLowerCase();
    var nodes = treeEl.querySelectorAll(".node");
    Array.prototype.forEach.call(nodes, function (n) {
      var hit = q === "" || (n.textContent || "").toLowerCase().indexOf(q) !== -1;
      n.classList.toggle("dim", q !== "" && !hit);
      n.classList.toggle("hit", q !== "" && hit);
    });
  }

  function setScale(value) {
    state.scale = Math.min(2, Math.max(0.4, value));
    treeEl.style.setProperty("--scale", String(state.scale));
  }

  async function load() {
    var clanVal = document.getElementById("clanSelect").value;
    var maxGen = parseInt(document.getElementById("maxGen").value, 10) || 5;
    var rootAttr = treeEl.getAttribute("data-root") || "0";
    var url = (clanVal && clanVal !== "0")
      ? "/api/clans/" + encodeURIComponent(clanVal) + "/tree?maxGen=" + maxGen
      : "/api/tree?root=" + encodeURIComponent(rootAttr) + "&maxGen=" + maxGen;

    try {
      var res = await fetch(url);
      if (!res.ok) throw new Error("HTTP " + res.status);
      var data = await res.json();
      renderTree(data.root);
    } catch (err) {
      treeEl.innerHTML = "";
      treeEl.appendChild(el("p", "muted", "Không tải được dữ liệu gia phả."));
    }
  }

  function addInfoRow(parent, label, value) {
    var p = el("p");
    p.appendChild(el("strong", null, label + ": "));
    p.appendChild(document.createTextNode(String(value)));
    parent.appendChild(p);
  }

  async function showDetail(id) {
    var popup = document.getElementById("popup");
    try {
      var res = await fetch("/api/persons/" + id);
      if (!res.ok) throw new Error("HTTP " + res.status);
      var data = await res.json();
      var person = data.person;

      document.getElementById("pName").textContent = person.full_name;

      var info = document.getElementById("pInfo");
      info.innerHTML = "";
      addInfoRow(info, "Giới tính", person.gender === 1 ? "Nam" : (person.gender === 2 ? "Nữ" : "Khác"));
      addInfoRow(info, "Năm sinh", person.birth_year || "Không rõ");
      addInfoRow(info, "Trạng thái", person.is_alive ? "Còn sống" : "Đã mất");

      if (data.spouses && data.spouses.length) {
        info.appendChild(el("p", null, "Vợ/Chồng:"));
        var ul = el("ul");
        data.spouses.forEach(function (sp) {
          var label = sp.partner_name + (sp.marriage_year ? " (" + sp.marriage_year + ")" : "");
          ul.appendChild(el("li", null, label));
        });
        info.appendChild(ul);
      }

      document.getElementById("pLink").setAttribute("href", "/persons/" + id);
      popup.hidden = false;
    } catch (err) {
      /* keep the popup closed on error */
    }
  }

  function closePopup() {
    var popup = document.getElementById("popup");
    if (popup) popup.hidden = true;
  }

  function init() {
    treeEl = document.getElementById("tree");
    if (!treeEl) return;

    document.getElementById("clanSelect").addEventListener("change", load);
    document.getElementById("maxGen").addEventListener("change", load);
    document.getElementById("searchBox").addEventListener("input", function (e) {
      state.query = e.target.value;
      applySearch();
    });
    document.getElementById("zoomIn").addEventListener("click", function () { setScale(state.scale + 0.1); });
    document.getElementById("zoomOut").addEventListener("click", function () { setScale(state.scale - 0.1); });

    var a3 = document.getElementById("a3");
    if (a3) {
      a3.addEventListener("change", function () {
        document.body.classList.toggle("print-a3", a3.checked);
      });
    }

    var popup = document.getElementById("popup");
    var popupClose = document.getElementById("popupClose");
    if (popupClose) popupClose.addEventListener("click", closePopup);
    if (popup) popup.addEventListener("click", function (e) { if (e.target === popup) closePopup(); });

    load();
  }

  document.addEventListener("DOMContentLoaded", init);
})();
