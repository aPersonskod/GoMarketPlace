export class ApiHelper {
    dev = (isHttps) => {
        return isHttps ? 'https://localhost:700' : 'http://localhost:808';
    };
    base = (serviceName) => {
        return `http://${serviceName}:8080`;
    };
    orchestratorServiceBaseAddress =       `${this.dev(false)}4/api`;
    buyServiceBaseAddress =                `${this.dev(false)}3/api/buy-service`;
    orderServiceBaseAddress =              `${this.dev(false)}2/api/order-service`;
    productServiceBaseAddress =            `${this.dev(false)}1/api/product-service`;
    userServiceBaseAddress =               `${this.dev(false)}0/api/user-service`;

    getUser = async () => {
        let token = this.getAccessToken();
        if (token === null) return null;
        let url = `${this.userServiceBaseAddress}/`;
        let options = {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
        };

        try{
            const response = await fetch(url, options);
            switch(response.status) {
                case 200:
                    const userDto = await response.json();
                    return userDto;

                case 401:
                    console.log("401 Unauthorized - Authentication failed");
                    //localStorage.removeItem('uToken');
                    return null;
                    
                case 403:
                    console.log("403 Forbidden - Insufficient permissions");
                    return null;
                default:
                    console.log(`Unhandled HTTP status: ${response.status}`);
                    return null;
            }
        }
        catch(e){
            console.log("Network or CORS error:", e);
            return null;
        }
    }

    getAccessToken = () => {
        return localStorage.getItem('uToken');
    }
}