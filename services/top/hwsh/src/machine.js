const endpoint = 'https://api.hello-world.sh/machine/';

export const machineStatus = async (token) => {
    await fetch(endpoint, {
        method: "GET",
        headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'Authorization': `Bearer: ${token}`
        }
    }).then((response) => {
        const result = response.json();
        alert(`machine status: ${result['status']}`)
    });
};

export const machineStart = async (token) => {
    await fetch(endpoint, {
        method: "PATCH",
        headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'Authorization': `Bearer: ${token}`
        }
    }).then((response) => {
        const result = response.json();
        alert(`machine: ${result['message']}, redirect: ${result['redirect_link']}`)
    });
};

export const machineStop = async (token) => {
    await fetch(endpoint, {
        method: "DELETE",
        headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'Authorization': `Bearer: ${token}`
        }
    }).then((response) => {
        const result = response.json();
        alert(`machine: ${result['message']}`)
    });
};