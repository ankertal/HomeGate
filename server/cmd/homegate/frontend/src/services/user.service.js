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

  createCommandJson(user, cmd) {
    var item = {}
    item[cmd] = true;
    item["gate_name"] = user.my_gate;
    return item
  }

  triggerCommand(user, cmd) {
    var postData = this.createCommandJson(user, cmd)
    return axios
      .post(API_URL + 'command', postData, { headers: authHeader() })
      .then(response => {
        if (response.data.status) {
          console.log('open returned' + response.data)
        }

        return response.data;
      });
  }

}

export default new UserService();
