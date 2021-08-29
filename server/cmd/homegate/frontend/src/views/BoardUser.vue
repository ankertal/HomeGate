<template>
  <div class="container">
    <header class="jumbotron">
      <h3>
        HomeGate's <strong>{{ currentUser.username }}</strong> Settings:
      </h3>
    </header>
    <div v-if="!successful">
      <div
        v-if="content"
        class="alert"
        :class="successful ? 'alert-success' : 'alert-danger'"
      >
        {{ content }}
      </div>
    </div>

    <div v-if="successful">
      <strong>Gates:</strong>
      <ul>
        <li v-for="(gate, index) in content.gates" :key="index">{{ gate }}</li>
      </ul>

      <div id="add-gate">
        <button v-on:click="addGate">Add Gate</button>
      </div>

      <table class="gates">
        <thead>
          <tr>
            <th>Name</th>
            <th>Owner</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in content.gates" :key="item">
            <td>{{ item }}</td>
            <td>
                <span v-if="item== content.my_gate"><font-awesome-icon icon="check" /></span>
                <span v-else><font-awesome-icon icon="times" /></span>
              </span>
            </td>
            <td>
              <button class="btn add">Add</button>
              <button class="btn edit">Edit</button>
            </td>
          </tr>
        </tbody>
      </table>
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
      message: "Hello Vue!",
      successful: false,
    };
  },
  mounted() {
    UserService.getUserBoard().then(
      (response) => {
        this.content = response.data;
        this.successful = true;
      },
      (error) => {
        this.successful = false;
        this.content =
          (error.response && error.response.data && error.response.data.message) ||
          error.message ||
          error.toString();
      }
    );
  },
  methods: {
    addGate() {
      this.content.gates.push("yaron-gate");
    },
  },
};
</script>
