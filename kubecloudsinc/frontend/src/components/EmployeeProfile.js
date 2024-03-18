import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';

function EmployeeProfile() {
  const { id } = useParams();
  const [employee, setEmployee] = useState(null);

  useEffect(() => {
    fetch(`/v2/employee/${id}`)
      .then(response => response.json())
      .then(data => setEmployee(data))
      .catch(error => console.error('Error:', error));
  }, [id]);

  return (
    <div>
      {employee ? (
        <div>
          <h1>{employee.name}</h1>
          {/* Display more employee details here */}
        </div>
      ) : (
        <p>Loading...</p>
      )}
    </div>
  );
}

export default EmployeeProfile;
