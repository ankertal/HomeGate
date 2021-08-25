<template>
  <div class="container">
    <header class="jumbotron">
      <h3>
        HomeGate's <strong>{{ currentUser.username }}</strong> Settings:
      </h3>
    </header>
    <strong>Gates:</strong>
    <ul>
      <li v-for="(gate, index) in content.gates" :key="index">{{ gate }}</li>
    </ul>
    <div id="counter">Counter: {{ counter }}</div>
    <div id="two-way-binding">
      <p>{{ message }}</p>
      <input v-model="message" />
    </div>
    <div id="event-handling">
      <p>{{ message }}</p>
      <button v-on:click="reverseMessage">Reverse Message</button>
    </div>
  </div>
</template>

<script>
import UserService from "../services/user.service";

export default {
  name: "User",
  computed: {
    currentUser() {
      return this.$store.state.auth.user;
    },
  },
  data() {
    return {
      content: "",
      counter: 0,
      message: "Hello Vue!",
    };
  },
  mounted() {
    UserService.getUserBoard().then(
      (response) => {
        this.content = response.data;
      },
      (error) => {
        this.content =
          (error.response && error.response.data && error.response.data.message) ||
          error.message ||
          error.toString();
      }
    );

    setInterval(() => {
      this.counter++;
    }, 1000);
  },
  methods: {
    reverseMessage() {
      this.message = this.message.split("").reverse().join("");
    },
  },
};
</script>
