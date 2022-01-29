CREATE TABLE employee
(
    id   INT PRIMARY KEY,
    name VARCHAR(50)
);

CREATE TABLE salary
(
    employee_id   INT PRIMARY KEY,
    salary_amount INT
);

ALTER TABLE salary
    ADD CONSTRAINT FK_Salary_Employee FOREIGN KEY (employee_id)
        REFERENCES employee (id);
