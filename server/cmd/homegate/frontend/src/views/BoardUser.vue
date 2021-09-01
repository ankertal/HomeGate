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
          <b-dropdown
            :variant="gateVariantInList(index)"
            class="m-2"
            size="md"
            :text="gate"
          >
          
              <b-dropdown-item @click="gateCommand(index, 'is_open')" active
                >Open</b-dropdown-item
              >
              <b-dropdown-item @click="gateCommand(index, 'is_close')"
                >Close</b-dropdown-item
              >
              <span v-if="isMyGate(index)">
                <b-dropdown-item disabled>Delete</b-dropdown-item>
              </span>
              <span v-else>
                <b-dropdown-item @click="deleteGate(index)">Delete</b-dropdown-item>
              </span>
            </b-dropdown>
          </b-dropdown>
        </li>
      </ul>

      <!-- <div id="gate-list">
        <form v-on:submit.prevent="addGate">
          <label for="new-gate">Add a gate to {{ currentUser.username }}: </label>
          &nbsp;&nbsp;
          <input v-model="newGateText" id="new-gate" placeholder="e.g. gate-XXXXXX" />
          <button>Add</button>
        </form>
      </div> -->
    </div>

    <div v-if="successful">
      <strong>My ({{ currentUser.my_gate }}) friends: </strong>
      <ul id="uses-list">
        <li v-for="(item, index) in this.content.users" :key="item">
          {{ item }}
        </li>
      </ul>

      <b-input-group prepend="Friend:" class="mb-3">
        <b-form-input v-model="email" :state="isEmailValid(this.email)"></b-form-input>
        <b-input-group-append>
          <b-button @click="addUser(currentUser.my_gate)" variant="outline-success"
            >Add</b-button
          >
          <b-button @click="deleteUser(currentUser.my_gate)" variant="info"
            >Delete</b-button
          >
        </b-input-group-append>
      </b-input-group>
    </div>

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
      msg: [],
      email: "",
      ismissCountDown: null,
      showDismissibleAlert: false,
      alertVariant: "info",
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
    gateVariantInList(index) {
      if (this.isMyGate(index)) {
        return 'outline-danger';
      }
      return 'primary';
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
    addUser(gateName) {
      if (/^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/.test(this.email)) {
        if (!this.content.users.includes(this.email)) {
          this.content.users.push(this.email);
        } else {
        }
        this.email = "";
      }
    },
    deleteUser(gateName) {
      if (this.email == this.currentUser.email) {
        return;
      }

      var index = this.content.users.indexOf(this.email);
      if (index > -1) {
        this.content.users.splice(index, 1);
      }
    },
    isEmailValid(value) {
      if (/^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/.test(value)) {
        return true;
      }
      return false;
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
