WITH updated_sla_condition
         AS (UPDATE cases.sla_condition SET updated_at = $1,updated_by = $2,name = $3 WHERE dc = $4 AND id = $5 RETURNING id,name,created_at,updated_at,reaction_time_hours,reaction_time_minutes,resolution_time_hours,resolution_time_minutes,sla_id,created_by,updated_by),
     insert_priorities AS (INSERT INTO cases.priority_sla_condition (created_at, updated_at, created_by, updated_by,
                                                                     sla_condition_id, priority_id,
                                                                     dc) SELECT $1, $1, $2, $2, 8::bigint, priority_id, 1
                                                                         FROM unnest(ARRAY [[1,2,3]]::bigint[]) AS priority_id ON CONFLICT (sla_condition_id,priority_id) DO NOTHING),
     update_priorities AS (UPDATE cases.priority_sla_condition SET updated_at = $1,updated_by = $2 WHERE
         sla_condition_id = 8::bigint AND priority_id IN (SELECT unnest(ARRAY [[1,2,3]]::bigint[]))),
     deleted_priorities AS (DELETE FROM cases.priority_sla_condition WHERE sla_condition_id = 8::bigint AND
                                                                           priority_id NOT IN
                                                                           (SELECT unnest(ARRAY [[1,2,3]]::bigint[])) RETURNING sla_condition_id,priority_id)
SELECT updated_sla_condition.id,
       updated_sla_condition.name,
       updated_sla_condition.created_at,
       updated_sla_condition.updated_at,
       updated_sla_condition.reaction_time_hours,
       updated_sla_condition.reaction_time_minutes,
       updated_sla_condition.resolution_time_hours,
       updated_sla_condition.resolution_time_minutes,
       updated_sla_condition.sla_id,
       updated_sla_condition.created_by                        AS created_by_id,
       COALESCE(c.name::text, c.username)                      AS created_by_name,
       updated_sla_condition.updated_by                        AS updated_by_id,
       COALESCE(u.name::text, u.username)                      AS updated_by_name,
       json_agg(json_build_object('id', p.id, 'name', p.name)) AS priorities
FROM updated_sla_condition
         LEFT JOIN directory.wbt_user u ON u.id = updated_sla_condition.updated_by
         LEFT JOIN directory.wbt_user c ON c.id = updated_sla_condition.created_by
         LEFT JOIN cases.priority_sla_condition psc ON updated_sla_condition.id = psc.sla_condition_id
         LEFT JOIN cases.priority p ON p.id = psc.priority_id
GROUP BY updated_sla_condition.id, updated_sla_condition.name, updated_sla_condition.created_at,
         updated_sla_condition.updated_at, updated_sla_condition.reaction_time_hours,
         updated_sla_condition.reaction_time_minutes, updated_sla_condition.resolution_time_hours,
         updated_sla_condition.resolution_time_minutes, updated_sla_condition.sla_id, updated_sla_condition.created_by,
         updated_sla_condition.updated_by, u.id, c.id;