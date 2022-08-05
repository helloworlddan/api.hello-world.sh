const endpoint = 'https://api.hello-world.sh/machine/';

export const machineStatus = (token) => {
    const response = await fetch(endpoint, {
        method: "GET",
        headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'Authorization': `Bearer: ${token}`
        }
    });

    if (!response.ok) {
        alert(`error: ${response.status}`);
    }

    const result = await response.json();

    alert(`machine status: ${result['status']}`)
};

export const machineStart = (token) => {
    const response = await fetch(endpoint, {
        method: "PATCH",
        headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'Authorization': `Bearer: ${token}`
        }
    });

    if (!response.ok) {
        alert(`error: ${response.status}`);
    }

    const result = await response.json();

    alert(`machine: ${result['message']}, redirect: ${result['redirect_link']}`)
};

export const machineStop = (token) => {
    const response = await fetch(endpoint, {
        method: "DELETE",
        headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'Authorization': `Bearer: ${token}`
        }
    });

    if (!response.ok) {
        alert(`error: ${response.status}`);
    }

    const result = await response.json();

    alert(`machine: ${result['message']}`)
};