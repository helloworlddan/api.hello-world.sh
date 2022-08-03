import React, { useContext } from "react";
import { UserContext } from "../providers/UserProvider";
import {auth} from "../firebase";
import {machineStatus, machineStart, machineStop} from "../machine";
const Dashboard = () => {
  const user = useContext(UserContext);
  const {photoURL, displayName, email, uid, _lat} = user;
  console.log(user);
  

  return (
    <div className = "mx-auto w-11/12 md:w-2/4 py-8 px-4 md:px-8">
      <div className="flex border flex-col items-center md:flex-row md:items-start border-blue-400 px-3 py-4">
        <div>
        <h3 className = "text-2xl font-semibold">Welcome, {displayName} ({email})!</h3>
        </div>
      </div>
      <button className = "w-full py-3 bg-yellow-600 mt-4 text-white" onClick = {() => {auth.signOut()}}>Sign out</button>
      
      <button className = "w-full py-3 bg-blue-600 mt-4 text-white" onClick = {() => {machineStatus()}}>Status</button>
      
      <button className = "w-full py-3 bg-green-600 mt-4 text-white" onClick = {() => {machineStart()}}>Start</button>

      <button className = "w-full py-3 bg-red-600 mt-4 text-white" onClick = {() => {machineStop()}}>Stop</button>

      <div className = "md:pl-4">Debug Info</div>
      <div>
        <table style={{
            tableLayout: "fixed"
            }}>
        <tbody style={{
            border: "1px solid black"
            }}>
          <tr style={{
            border: "1px solid black"
            }}>
            <td><code>{uid}</code></td>
          </tr >
          <tr style={{
            border: "1px solid black"
            }}>
            <td style={{
            wordBreak: "break-all"
            }}><code>{_lat}</code></td>
          </tr>
        </tbody>
        </table>
      </div>
    </div>
  ) 
};

export default Dashboard;