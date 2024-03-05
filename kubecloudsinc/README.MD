# **KubeCloudsInc Employee Management API**

The Employee Management API provides a set of endpoints for managing employee records. This API includes functionalities for user authentication, retrieving employee details, adding new employees, and updating existing employee information.

## **Endpoints Overview:**

### **Generate Login Token**
**Endpoint:** /v2/login

**Method:** POST

**Description:** Authenticates users and provides a token for accessing protected endpoints. This endpoint does not require pre-existing authorization but returns credentials needed for further API interactions.

### **Get Employees**
**Endpoint:** /v2/employees

**Method:** GET

**Authorization Required:** admin, editor, viewer

**Description:** Retrieves a list of all employees. Accessible by users with admin, editor, or viewer roles.

### **Get Employee**
**Endpoint:** /v2/employee

**Method:** GET

**Authorization Required:** admin, editor, viewer

**Description:** Retrieves detailed information for a specific employee. Users need to provide the employee's ID as a query parameter. Accessible by users with admin, editor, or viewer roles.

### **Add Employee**
**Endpoint:** /v2/employee

**Method:** POST

**Authorization Required:** admin, editor

**Description:** Adds a new employee record to the database. Requires sending employee details in the request body. Accessible by users with admin or editor roles.

### **Update Employee**
**Endpoint:** /v2/employee/{employeeId}

**Method:** PUT

**Authorization Required:** admin, editor

**Description:** Updates details for an existing employee. The employee's ID is specified in the URL, and the details to be updated are sent in the request body. Accessible by users with admin or editor roles.

### **Delete Employee**
**Endpoint:** /v2/employee/{employeeId}

**Method:** DELETE

**Authorization Required:** admin

**Description:** Deletes details for an existing employee. The employee's ID is specified in the URL. Accessible by users with admin roles.

### **Authorization**
Access to most endpoints requires authorization. After logging in, users will receive a token which must be included in the Authorization header of subsequent requests. The role associated with the user's token determines which endpoints can be accessed.

### **Request & Response Formats**
**Login Request:** Send credentials in the request body in JSON format.

**Protected Endpoints:** Include the Authorization header with the token received upon login.

**GET Requests:** For endpoints that fetch data, parameters (if any) should be included in the query string.

**POST/PUT Requests:** Send data in the request body in JSON format. Refer to the API documentation for detailed schema definitions.

### **Error Handling**
Responses to unsuccessful requests will include an appropriate HTTP status code and a JSON object containing an error message. Clients should handle these responses gracefully.

### **Docker**
docker build -t kube .

docker run -d -p 8080:8080 -e DATABASE_DSN="admin/Jaffa123@10.10.12.130:1521/GHGWE1" kube   

Note: `admin/Jaffa123@10.10.12.130:1521/GHGWE1` is a dummy DSN, please replace it with actual one while running/testing