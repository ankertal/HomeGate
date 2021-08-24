import axios from 'axios';
import authHeader from './auth-header';

const API_URL = 'http://localhost/static/';

class UserService {
  getPublicContent() {
    return axios.get(API_URL + "all.html");
  }

  getUserBoard() {
    return axios.get(API_URL + 'user.html', { headers: authHeader() });
  }

  getModeratorBoard() {
    return axios.get(API_URL + 'mod', { headers: authHeader() });
  }

  getAdminBoard() {
    return axios.get(API_URL + 'admin', { headers: authHeader() });
  }
}

export default new UserService();