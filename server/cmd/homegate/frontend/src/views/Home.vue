<template>
  <div class="container">
    <header class="jumbotron">
      <div v-if="!successful">
        <h3>HomeGate:</h3>
        <h4
          v-if="content"
          class="alert"
          :class="successful ? 'alert-success' : 'alert-danger'"
        >
          {{ content }}
        </h4>
      </div>
      <div v-if="successful">
        <h3>{{ content }}</h3>
      </div>
    </header>
    <h5>Created by:</h5>
    <br />
    <div>
      <p>Yaron Weinsberg</p>
      <img
        class="profile-img-card"
        src="@/assets/yaron.jpeg"
        alt="Yaron Weinsberg"
        width="100"
        height="100"
        hspace="50"
      />
    </div>
    <br />
    <br />
    <div>
      <p>Tal Anker</p>
      <img
        class="profile-img-card"
        src="@/assets/anker.jpeg"
        alt="Tal Anker"
        width="100"
        height="100"
        hspace="50"
      />
    </div>
  </div>
</template>

<script>
import UserService from "../services/user.service";

export default {
  name: "Home",
  data() {
    return {
      content: "",
      successful: false,
    };
  },
  mounted() {
    UserService.getPublicContent().then(
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
};
</script>

<style scoped>
.profile-img-card {
  width: 96px;
  height: 96px;
  margin-left: 30px;
  display: block;
  -moz-border-radius: 50%;
  -webkit-border-radius: 50%;
  border-radius: 50%;
}
</style>
