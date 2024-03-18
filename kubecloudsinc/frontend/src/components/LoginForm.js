import React, { useState } from 'react';
import '../App.css'; // Make sure the path is correct for your project structure

function LoginForm({ onLoginSuccess }) {
  const [credentials, setCredentials] = useState({ username: '', password: '' });
  const [error, setError] = useState('');

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setCredentials(prevState => ({ ...prevState, [name]: value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      // Here's where you make the call to your API endpoint for login
      const response = await fetch('http://192.168.1.31:8080/v2/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json'},
        body: JSON.stringify(credentials),
      });

      // Check if the response from the server is ok (status code 200)
      if (!response.ok) {
        // If the server response is not OK, extract the error message from the response
        const errorData = await response.json();
        throw new Error(errorData.message || 'Login failed');
      }

      // If the response is ok, extract the token from the response data
      const data = await response.json();
      localStorage.setItem('token', data.token); // Store the token in local storage
      onLoginSuccess(data.token); // Call the onLoginSuccess callback with the token

    } catch (err) {
      // If there's an error, set the error state with the error message
      setError(err.message);
    }
  };

  return (
    <div className="login-container">
      <form className="login-form" onSubmit={handleSubmit}>
        <input
          className="login-input"
          id="username"
          name="username"
          type="text"
          placeholder="Username"
          value={credentials.username}
          onChange={handleInputChange}
          required
        />
        <input
          className="login-input"
          id="password"
          name="password"
          type="password"
          placeholder="Password"
          value={credentials.password}
          onChange={handleInputChange}
          required
        />
        <button className="login-button" type="submit">Login</button>
        {error && <div className="error">{error}</div>}
      </form>
    </div>
  );
}

export default LoginForm;




/*import React, { useState } from 'react';

function LoginForm({ onLoginSuccess }) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const handleLogin = async (e) => {
    e.preventDefault();
    // Post request to backend
    try {
      const response = await fetch('http://localhost:8080/v2/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
      });
      const data = await response.json();
      if (response.ok) {
        localStorage.setItem('token', data.token); // Adjust according to your response structure
        onLoginSuccess(data.token);
      } else {
        throw new Error(data.message || 'Login failed');
      }
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div>
      <form onSubmit={handleLogin}>
        <input
          type="text"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          placeholder="Username"
          required
        />
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder="Password"
          required
        />
        <button type="submit">Login</button>
        {error && <p>{error}</p>}
      </form>
    </div>
  );
}

export default LoginForm;
*/