import React, { useState, useEffect } from 'react';

function EmployeeList() {
  const [employees, setEmployees] = useState([]);
  const [error, setError] = useState('');

  useEffect(() => {
    const token = localStorage.getItem('token'); // Get the token from local storage

    if (!token) {
      setError('No authentication token. Please log in.');
      return;
    }

    fetch('http://192.168.1.31:8080/v2/employees', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`, // Use the token for authorization
        'Content-Type': 'application/json',
      },
    })
    .then(response => {
      if (!response.ok) {
        // If response is not OK, throw an error
        throw new Error(`HTTP status ${response.status}`);
      }
      return response.json();
    })
    .then(data => {
      setEmployees(data); // Set the employees data on state
    })
    .catch(error => {
      console.error('Failed to fetch:', error);
      setError('Failed to fetch employees.'); // Set an error message on state
    });
  }, []);

  return (
    <div>
      <h1>Employee List</h1>
      {error && <div className="error">{error}</div>}
      <ul>
        {employees.map(employee => (
          <li key={employee.employeeId}>
            {employee.firstName} {employee.lastName}
          </li>
        ))}
      </ul>
    </div>
  );
}

export default EmployeeList;

/*import React, { useState, useEffect } from 'react';

function EmployeeList() {
  const [employees, setEmployees] = useState([]);
  const [error, setError] = useState(null);

  useEffect(() => {
    const token = ''; // Replace with your actual token from authentication
    fetch('http://localhost:8080/v2/employees', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}` // Include the authorization header
      },
    })
    .then(response => {
      if (!response.ok) {
        throw new Error(`HTTP status ${response.status}`);
      }
      return response.json();
    })
    .then(data => setEmployees(data))
    .catch(error => {
      console.error('Fetch error:', error);
      setError(error.message);
    });
  }, []);

  return (
    <div>
      <h1>Employee List</h1>
      {error && <p>Error: {error}</p>}
      <ul>
        {employees.map(employee => (
          <li key={employee.employeeId}>{employee.firstName} {employee.lastName}</li>
        ))}
      </ul>
    </div>
  );
}

export default EmployeeList;
*/

// import React, { useState, useEffect } from 'react';
// import { Link } from 'react-router-dom';

// function EmployeeList() {
//   const [employees, setEmployees] = useState([]);

//   useEffect(() => {
//     fetch('/v2/employees')
//       .then(response => response.json())
//       .then(data => setEmployees(data))
//       .catch(error => console.error('Error:', error));
//   }, []);

//   return (
//     <div>
//       <h1>Employee List</h1>
//       <ul>
//         {employees.map(employee => (
//           <li key={employee.id}>
//             {employee.name} - <Link to={`/employee/${employee.id}`}>View Profile</Link>
//           </li>
//         ))}
//       </ul>
//     </div>
//   );
// }

// export default EmployeeList;
