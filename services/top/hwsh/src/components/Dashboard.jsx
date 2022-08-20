import React, { useContext, useEffect, useState, useRef } from "react";
import { Navigate } from 'react-router-dom'
import { UserContext } from "../providers/UserProvider";
import { auth } from "../firebase";

const Dashboard = () => {
  const user = useContext(UserContext);
  const { _lat } = user;
  const [machine, setMachine] = useState();
  const endpoint = 'https://api.hello-world.sh/machine/';

  const executeRequest = async (requestMethod, token) => {
    const response = await fetch(endpoint, {
      method: requestMethod,
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        'Authorization': `Bearer ${token}`
      }
    })
    if (!response.ok) {
      console.log(response);
    }
    return await response.json();
  }

  const machineStatus = async (token) => {
    executeRequest("GET", token).then((result) => {
      setMachine(result);
    });
  };

  const machineStart = async (token) => {
    executeRequest("PATCH", token).then((result) => {
      setMachine(result);
    });
  };

  const machineStop = async (token) => {
    executeRequest("DELETE", token).then((result) => {
      setMachine(result);
    });
  };

  useEffect(() => {
    setMachine({
      "status": "refreshing ...",
      "redirect_link": "none"
    });
    machineStatus(_lat);
    if (machine && machine["status"] && machine["status"] === "RUNNING") {
      return <Navigate to={machine["redirect_link"]} />
    }
  }, []);

  return (
    <div className="mx-auto w-11/12 md:w-2/4 py-8 px-4 md:px-8">
      <div className="flex border flex-col items-center md:flex-row md:items-start border-blue-400 px-3 py-4">
        {machine ? (
          <div>
            <p>Status: {machine["status"] || "-"}</p>
            <p><a href={machine["redirect_link"] || "https://remotedesktop.google.com"}>Connect</a></p>
          </div>
        ) : (
          <div>
            <p>retrieving....</p>
          </div>
        )
        }
      </div>
      <button className="w-full py-3 bg-yellow-600 mt-4 text-white" onClick={() => { auth.signOut() }}>Sign out</button>

      <button className="w-full py-3 bg-blue-600 mt-4 text-white" onClick={() => { machineStatus(_lat) }}>Status</button>

      <button className="w-full py-3 bg-green-600 mt-4 text-white" onClick={() => { machineStart(_lat) }}>Start</button>

      <button className="w-full py-3 bg-red-600 mt-4 text-white" onClick={() => { machineStop(_lat) }}>Stop</button>
    </div>
  )
};

export default Dashboard;