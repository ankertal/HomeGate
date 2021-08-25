<template>
  <div class="container">
    <header class="jumbotron">
      <h3>{{content}}</h3>
    </header>
      <h5>Created by:</h5>
      <br/>
      <div> 
        <p> Yaron Weinsberg </p> 
        <img src="@/assets/yaron.jpeg" alt="Girl in a jacket" width="100" height="100" hspace="50"/>
      </div>
      <br/>
      <br/>
      <div> 
        <p> Tal Anker </p> 
        <img src="@/assets/anker.jpeg" alt="Girl in a jacket" width="100" height="100" hspace="50"/>
      </div>
  </div>
</template>

<script>
import UserService from '../services/user.service';

export default {
  name: 'Home',
  data() {
    return {
      content: ''
    };
  },
  mounted() {
    UserService.getPublicContent().then(
      response => {
        this.content = response.data;
      },
      error => {
        this.content =
          (error.response && error.response.data && error.response.data.message) ||
          error.message ||
          error.toString();
      }
    );
  }
};
</script>
