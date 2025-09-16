-- Clean up the table from deleted assignee_id values,
-- i.e., those referencing deleted contacts
UPDATE cases.service_catalog sc
SET assignee_id = NULL
WHERE assignee_id IS NOT NULL
  AND NOT EXISTS (
      SELECT 1
      FROM contacts.contact c
      WHERE c.id = sc.assignee_id
  );

-- Add a foreign key
ALTER TABLE cases.service_catalog
ADD CONSTRAINT fk_service_catalog_assignee
FOREIGN KEY (assignee_id)
REFERENCES contacts.contact(id)
ON DELETE SET NULL;