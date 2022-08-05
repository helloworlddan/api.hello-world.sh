import { useContext } from "react";
import { UserContext } from "../providers/UserProvider";

const endpoint = 'https://api.hello-world.sh/machine/';

export const machineStatus = () => {
    const user = useContext(UserContext);
    const {_lat} = user;

    const response = fetch(endpoint, {
        method: "GET",
        headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'Authorization': `Bearer: ${_lat}`
        }
    });

    if (!response.ok) {
        alert(`error: ${response.status}`);
    }

    const result = response.json();

    alert(`machine status: ${result['status']}`)
};

export const machineStart = () => {
    const user = useContext(UserContext);
    const {_lat} = user;

    const response = fetch(endpoint, {
        method: "PATCH",
        headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'Authorization': `Bearer: ${_lat}`
        }
    });

    if (!response.ok) {
        alert(`error: ${response.status}`);
    }

    const result = response.json();

    alert(`machine: ${result['message']}, redirect: ${message['redirect_link']}`)
};

export const machineStop = () => {
    const user = useContext(UserContext);
    const {_lat} = user;

    const response = fetch(endpoint, {
        method: "DELETE",
        headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'Authorization': `Bearer: ${_lat}`
        }
    });

    if (!response.ok) {
        alert(`error: ${response.status}`);
    }

    const result = response.json();

    alert(`machine: ${result['message']}`)
};