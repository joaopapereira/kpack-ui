import axios from 'axios';
class HttpApi {
    async request(method, url, data) {
        try {
            const response = await axios({
                method: method,
                url: url,
                data: data,
                headers: {
                    'content-type': 'application/json',
                }
            });

            return response.data;
        } catch (err) {
            const error = err;
            console.error(error.response.status);
            throw err;
        }
    }
}
const httpApi = new HttpApi();
export default httpApi;