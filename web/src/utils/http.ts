import axios from 'axios';

export interface HttpAPI {
  request(method, url, data): Promise<object>;
}

class HttpApiImpl implements HttpAPI {
  async request(method, url, data): Promise<object> {
    try {
      const response = await axios({
        method: method,
        url: url,
        data: data,
        headers: {
          'content-type': 'application/json'
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

const httpApi = new HttpApiImpl();
export default httpApi;
