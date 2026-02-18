-- TREE FROM ROOT
WITH RECURSIVE family_tree AS (
    SELECT id, full_name, gender, father_id, mother_id, is_alive, clan_id, 0 AS depth
    FROM person
    WHERE id = $1
    UNION ALL
    SELECT p.id, p.full_name, p.gender, p.father_id, p.mother_id, p.is_alive, p.clan_id, ft.depth + 1
    FROM person p
    JOIN family_tree ft
      ON p.father_id = ft.id OR p.mother_id = ft.id
)
SELECT * FROM family_tree
ORDER BY depth, id;

