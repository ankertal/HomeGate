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
                  <b-button
                    pill
                    variant="outline-success"
                    size="sm"
                    @click="gateCommand(index, 'is_open')"
                    >Open</b-button
                  >
                </li>
                &nbsp;&nbsp;
                <li class="nav-item">
                  <b-button
                    pill
                    variant="outline-warning"
                    size="sm"
                    @click="gateCommand(index, 'is_close')"
                    >Close</b-button
                  >
                </li>
                &nbsp;&nbsp;
                <li class="nav-item active">
                  <b-button
                    pill
                    variant="outline-danger"
                    size="sm"
                    @click="deleteGate(index)"
                    :disabled="isMyGate(index)"
                  >
                    Delete
                  </b-button>
                </li>
                &nbsp;&nbsp;&nbsp;&nbsp;
                <li v-if="isMyGate(index)" class="nav-item active">
                  <b-button pill variant="outline-dark" size="sm" @click="addUser(index)"
                    >Add User</b-button
                  >&nbsp;
                  <input
                    class="text-left"
                    size="sm"
                    type="text"
                    v-model="email"
                    required
                  />
                  <span v-if="msg.email">{{ msg.email }}</span>
                  <!-- <div>
                    <b-dropdown split text="Options" class="m-2">
                      <b-dropdown-item @click="addGate()">Add User</b-dropdown-item>
                      <b-dropdown-item>Another action</b-dropdown-item>
                      <b-dropdown-item>Something else here...</b-dropdown-item>
                    </b-dropdown>
                  </div> -->
                </li>
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

    <br /><br /><br />
    <div v-if="successful">
      <strong>My ({{ currentUser.my_gate }}) friends: </strong>
      <ul id="uses-list">
        <li v-for="(item, index) in this.content.users">
          {{ item }}
        </li>
      </ul>
    </div>

    <br /><br /><br />
    <div v-if="this.showCommandStatus">
      <b-alert
        :show="dismissCountDown"
        :variant="alertVariant"
        dismissible
        @dismiss-count-down="countDownChanged"
      >
        {{ this.message }}</b-alert
      >
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
      message: "",
      cmdError: false,
      showCommandStatus: false,
      successful: false,
      newGateText: "",
      newUserText: "",
      msg: [],
      email: "",
      ismissCountDown: null,
      showDismissibleAlert: false,
      alertVariant: "info",
    };
  },
  watch: {
    email(value) {
      // binding this to the data value in the email input
      this.email = value;
      this.validateEmail(value);
    },
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
      if (this.newGateText != "") {
        this.content.gates.push(this.newGateText);
        this.newGateText = "";
      }
    },
    deleteGate(index) {
      const currentGate = this.content.gates[index];
      const myGate = this.content.user_gate;
      if (currentGate != myGate) {
        this.content.gates.splice(index, 1);
      }
    },
    isMyGate(index) {
      const currentGate = this.content.gates[index];
      const myGate = this.content.my_gate;
      return currentGate === myGate;
    },
    gateCommand(index, cmd) {
      this.message = "";
      this.cmdError = false;
      const currentGate = this.content.gates[index];
      // const myGate = this.content.my_gate;
      // if (currentGate == myGate) {
      // open the gate
      UserService.triggerCommand(this.currentUser, currentGate, cmd).then(
        (response) => {
          if (response.isAxiosError || response.is_error) {
            this.message = response.message;
            this.showCommandStatus = true;
            this.showAlert("danger");
          } else {
            this.cmdError = response.is_error;
            this.message = response.message;
            this.showCommandStatus = !this.cmdError;
            this.showAlert("info");
          }
        },
        (error) => {
          this.cmdError = true;
          this.message =
            (error.response && error.response.data && error.response.data.message) ||
            error.message ||
            error.toString();
          this.showAlert("danger");
        }
      );
      //}
    },
    addUser: function () {
      if (/^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/.test(this.email)) {
        this.content.users.push(this.email);
        this.newUserText = "";
      }
    },
    validateEmail(value) {
      if (/^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/.test(value)) {
        this.msg["email"] = "";
      } else {
        this.msg["email"] = "Invalid Email Address";
      }
    },
    countDownChanged(dismissCountDown) {
      this.dismissCountDown = dismissCountDown;
    },
    showAlert(variant) {
      this.dismissCountDown = 2;
      this.alertVariant = variant;
    },
  },
};
</script>
