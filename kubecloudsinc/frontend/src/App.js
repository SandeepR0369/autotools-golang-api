// App.js

// import React, { useState } from 'react';
// import LoginForm from './components/LoginForm';
// import EmployeeList from './components/EmployeeList';
// import AdminPanel from './components/AdminPanel'; // Import the new component
// import './App.css'; // Make sure to import the CSS file

// function App() {
//   const [token, setToken] = useState(localStorage.getItem('token'));

//   const handleLoginSuccess = (newToken) => {
//     localStorage.setItem('token', newToken); // Save the token in local storage
//     setToken(newToken); // Update state with the new token
//   };

//   return (
//     <div className="App">
//       <div className="header">
//         <h1 className="logo">KubeCloudsInc</h1>
//       </div>
//       {!token ? (
//         <LoginForm onLoginSuccess={handleLoginSuccess} />
//       ) : (
//         <>
//           <EmployeeList />
//           <AdminPanel />
//           {/* You can add additional components that should render after login here */}
//         </>
//       )}
//     </div>
//   );
// }

// export default App;

import React, { useState } from 'react';
import LoginForm from './components/LoginForm';
import EmployeeList from './components/EmployeeList';
//import AdminPanel from './components/AdminPanel'; 
import './App.css'; // Make sure to import the CSS file

function App() {
  const [token, setToken] = useState(localStorage.getItem('token'));

  const handleLoginSuccess = (token) => {
    setToken(token);
    localStorage.setItem('token', token); // Save the token in local storage
  };

  return (
    <div className="App">
      <div className="header">
        <h1 className="logo">KubeCloudsInc</h1>
      </div>
      {token ? (
        <EmployeeList />
      ) : (
        <LoginForm onLoginSuccess={handleLoginSuccess} />
      )}
    </div>
  );
}

export default App;

/*import React from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import EmployeeList from './components/EmployeeList';
import EmployeeProfile from './components/EmployeeProfile';
import './App.css';

function App() {
  return (
    <Router>
      <div>
        <Switch>
          <Route path="/" exact component={EmployeeList} />
          <Route path="/employee/:id" component={EmployeeProfile} />
        </Switch>
      </div>
    </Router>
  );
}

export default App;
*/
