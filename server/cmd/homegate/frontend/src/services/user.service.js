import axios from 'axios';
import authHeader from './auth-header';

const API_URL = 'http://localhost:80/';


class UserService {
  getPublicContent() {
    return axios.get(API_URL + "static/" + "all.html");
  }

  getUserBoard() {
    return axios.get(API_URL + 'user', { headers: authHeader() });
  }

  createCommandJson(user, gateName, cmd) {
    var item = {}
    item[cmd] = true;
    //item["gate_name"] = user.my_gate;
    item["gate_name"] = gateName;
    return item
  }

  triggerCommand(user, gateName, cmd) {
    var postData = this.createCommandJson(user, gateName, cmd)
    return axios
      .post(API_URL + 'command', postData, { headers: authHeader() })
      .then(response => {
        if (response.data.status) {
          console.log('triggerCommand returned' + response.data)
        }
        return response.data;
      })
      .catch(function (error) {
        if (error.response) {
          console.log(error.response.data);
          console.log(error.response.status);
          console.log(error.response.headers);
          return error.response.data;
        }
        return error
      });
  }

}

export default new UserService();
