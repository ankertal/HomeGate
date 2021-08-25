<template>
  <div class="container">
    <header class="jumbotron">
      <h3>
        HomeGate's <strong>{{currentUser.username}}</strong> Settings:
      </h3>
    </header>
    <p>
       {{content}}
    </p>
  </div>
</template>

<script>
import UserService from '../services/user.service';

export default {
  name: 'User',
  computed: {
    currentUser() {
      return this.$store.state.auth.user;
    }
  },
  data() {
    return {
      content: ''
    };
  },
  mounted() {
    UserService.getUserBoard().then(
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
