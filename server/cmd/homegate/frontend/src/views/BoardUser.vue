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
        <li v-for="(gate, index) in content.gates" :key="index">
          <nav class="navbar navbar-expand-lg navbar-light bg-light">
            <a class="navbar-brand" href="#"
              ><span v-if="isMyGate(index)" style="color: red">{{ gate }}</span>
              <span v-if="!isMyGate(index)" style="color: blue">{{ gate }}</span></a
            >
            <button
              class="navbar-toggler"
              type="button"
              data-toggle="collapse"
              data-target="#navbarNav"
              aria-controls="navbarNav"
              aria-expanded="false"
              aria-label="Toggle navigation"
            >
              <span class="navbar-toggler-icon"></span>
            </button>
            <div class="collapse navbar-collapse" id="navbarNav">
              <ul class="navbar-nav">
                <li class="nav-item">
                  <button @click="openGate(index)">Open</button>
                </li>
                <li class="nav-item">
                  <button @click="closeGate(index)">Close</button>
                </li>
                <li class="nav-item active">
                  <button @click="deleteGate(index)" :disabled="isMyGate(index)">
                    Delete
                  </button>
                </li>
                <!-- <li class="nav-item">
              <a class="nav-link disabled" href="#">Disabled</a>
            </li> -->
              </ul>
            </div>
          </nav>
        </li>
      </ul>

      <div id="gate-list">
        <form v-on:submit.prevent="addGate">
          <label for="new-gate">Add a gate to {{ currentUser.username }}: </label>
          &nbsp;&nbsp;
          <input v-model="newGateText" id="new-gate" placeholder="e.g. gate-XXXXXX" />
          <button>Add</button>
        </form>
      </div>
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
      newGateText: "",
    };
  },
  mounted() {
    UserService.getUserBoard().then(
      (response) => {
        this.content = response.data;
        this.successful = true;
        this.items = this.content.gates;
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
    addGate: function () {
      this.content.gates.push(this.newGateText);
      this.newGateText = "";
    },
    deleteGate(index) {
      const currentGate = this.content.gates[index];
      const myGate = this.content.my_gate;
      if (currentGate.trim() != myGate.trim()) {
        this.content.gates.splice(index, 1);
      }
    },
    isMyGate(index) {
      const currentGate = this.content.gates[index];
      const myGate = this.content.my_gate;
      return currentGate.trim() === myGate.trim();
    },
  },
};
</script>
