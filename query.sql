SELECT
    e.UserName UserName,
    m.UserName AS ParentUserName
FROM
    USER e
LEFT JOIN USER m ON  e.Parent =  m.ID;	