import React, { useContext } from "react";
import { Router } from "@reach/router";
import SignIn from "./SignIn";
import Dashboard from "./Dashboard";
import { UserContext } from "../providers/UserProvider";

function Application() {
  const user = useContext(UserContext);
  return (
        user ?
        <Dashboard />
      :
        <Router>
          <SignIn path="/" />
        </Router>
      
  );
}

export default Application;