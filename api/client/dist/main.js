
const apiGatewayUrl = "https://1lpy7pswjk.execute-api.us-west-2.amazonaws.com/dev";
const apiKey = "U332r4z0dV60z8gaVrNJvCOA0YQcmMBayQuudrg8";

const getServerList = async () => {
    const listEndpoint = `${apiGatewayUrl}/list`;
    const options = {
        method: 'GET',
        headers: {
            'x-api-key': apiKey,
        },
    };
    const serverList = await fetch(
        listEndpoint,
        options,
    );
    return serverList.json();
}

const createSteamConnectElement = (ipAddress) => {
    const steamConnectUrl = `steam://connect/${ipAddress}:27015`;
    const buttonElement = document.createElement('button');
    buttonElement.innerHTML = "Connect";
    buttonElement.setAttribute("type", "button");
    buttonElement.setAttribute("onclick", `window.open("${steamConnectUrl}");`);
    if(ipAddress == "") {
        buttonElement.setAttribute("disabled", true);
    }
    return buttonElement.outerHTML;
}

const startServer = async (acsHost) => {
    const startEndpoint = `${apiGatewayUrl}/start`;
    const data = {'acs-host': acsHost};
    const options = {
        method: 'POST',
        headers: {
            'x-api-key': apiKey,
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
    };
    console.log("Starting Server: ", acsHost);
    const startRequest = await fetch(startEndpoint, options);
}

const stopServer = async (acsHost) => {
    const stopEndpoint = `${apiGatewayUrl}/stop`;
    const data = {'acs-host': acsHost};
    const options = {
        method: 'POST',
        headers: {
            'x-api-key': apiKey,
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
    };
    console.log("Stopping Server: ", acsHost);
    const stopRequest = await fetch(stopEndpoint, options);
}

const createChangeStatusButton = (acsHost, status) => {
    const buttonElement = document.createElement('button');
    buttonElement.setAttribute("type", "button");
    switch(status) {
        case 'running':
            buttonElement.innerHTML = "Stop Server";
            buttonElement.setAttribute("onclick",`this.toggleAttribute("disabled");stopServer("${acsHost}");`);
            break;
        case 'stopped':
            buttonElement.innerHTML = "Start Server";
            buttonElement.setAttribute("onclick",`this.toggleAttribute("disabled");startServer("${acsHost}");`);
            break;
        default:
            buttonElement.innerHTML = "Stop Server";
            buttonElement.setAttribute("disabled", true);
    }
    return buttonElement.outerHTML;
}

const populateServerList = async () => {
    const serverListPromise = getServerList();
    const tableElement = document.createElement("table");
    const listContainer = document.getElementById("server-list");
    const serverList = await serverListPromise;
    let rowElement;
    let tdTemp;

    tableHeadElement = document.createElement("thead");
    rowElement = document.createElement("tr");
    for(const key in serverList[0]){
        tdTemp = document.createElement("th");
        switch(key) {
            case "InstanceID":
                tdTemp.innerHTML = "EC2 Instance ID";
                break;
            case "Name":
                tdTemp.innerHTML = "Map Name";
                break;
            case "PublicIPAddress":
                tdTemp.innerHTML = "Steam Connect";
                break;
            case "ActivePlayers":
                tdTemp.innerHTML = "Active Player Count";
                break;
            default:
                tdTemp.innerHTML = key;
        }
        rowElement.appendChild(tdTemp);
    }
    tdTemp = document.createElement("th");
    tdTemp.innerHTML = "Start/Stop";
    rowElement.appendChild(tdTemp);
    tableHeadElement.appendChild(rowElement);
    tableElement.appendChild(tableHeadElement);

    tableBodyElement = document.createElement("tbody");
    serverList.forEach( server => {
        rowElement = document.createElement("tr");
        for(const key in server){
            tdTemp = document.createElement("td");
            switch(key) {
                case "PublicIPAddress":
                    tdTemp.innerHTML = createSteamConnectElement(server[key]);
                    break;
                default:
                    tdTemp.innerHTML = server[key];
            }
            rowElement.appendChild(tdTemp);
        }
        tdTemp = document.createElement("td");
        tdTemp.innerHTML = createChangeStatusButton(server.Name, server.Status);
        rowElement.appendChild(tdTemp);
        tableBodyElement.appendChild(rowElement);
    } );
    tableElement.appendChild(tableBodyElement);
    listContainer.innerHTML = tableElement.outerHTML;
};

populateServerList();
setInterval(populateServerList, 30000);
