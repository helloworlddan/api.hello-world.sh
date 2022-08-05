const endpoint = 'https://api.hello-world.sh/machine/';

const executeRequest = async (requestMethod, token) => {
    const response = await fetch(endpoint, {
        credentials: 'include',
        method: requestMethod,
        headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'Authorization': `Bearer: ${token}`
        }
    })

    if (!response.ok) {
        console.log(response);
    }
    
    return await response.json();
}

export const machineStatus = async (token) => {
    executeRequest("GET", token).then((result) => {
        alert(`machine status: ${result['status']}`)
    });
};

export const machineStart = async (token) => {
    executeRequest("PATCH", token).then((result) => {
        alert(`machine: ${result['message']}, redirect: ${result['redirect_link']}`)
    });
};

export const machineStop = async (token) => {
    executeRequest("DELETE", token).then((response) => {
        const result = response.json();
        alert(`machine: ${result['message']}`)
    });
};
