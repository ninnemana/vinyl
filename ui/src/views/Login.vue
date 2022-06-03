<template>
  <v-content>
    <v-container fluid>
      <v-row>
        <v-col class='login' cols='8'>
          <v-row align='center' justify='center'>
            <v-card class='elevation-12'>
              <v-form @submit.prevent='handleSubmit' id='login-form' ref='form'>
                <v-card-text>
                  <v-text-field 
                    label='Login' 
                    name='login' 
                    prepend-icon='person' 
                    type='text' 
                    v-model='email'
                  />

                  <v-text-field
                    id='password'
                    label='Password'
                    name='password'
                    prepend-icon='lock'
                    type='password'
                    v-model='password'
                  />
                </v-card-text>
                <v-card-actions>
                  <v-spacer />
                  <v-btn type='submit' color='primary' class>Login</v-btn>
                </v-card-actions>
              </v-form>
            </v-card>
          </v-row>
          <!-- <v-row align='center' justify='center' class='social'>
            <v-card class='elevation-12'>
              <a href="/auth/github">
                <button class="btn-auth btn-github">
                  <span>Sign in with <b>GitHub</b></span>
                  <v-icon>fab fa-github</v-icon>
                </button>
              </a>
            </v-card>
          </v-row> -->
          <v-row align='center' justify='center' class='create-account'>
            <router-link to='/account/create'>Create Account</router-link>
            <span>|</span>
            <router-link to='/account/reset-password'>Reset Password</router-link>
          </v-row>
        </v-col>
      </v-row>
    </v-container>
    <v-snackbar v-model='ShowError' right>
      {{ LoginError }}
      <v-btn color="pink" text>
        Close
      </v-btn>
    </v-snackbar>
  </v-content>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
// import { mapGetters } from 'vuex'
import CreateAccount from './CreateAccount.vue';
import ResetPassword from './ResetPassword.vue';

@Component({
	components: {
		CreateAccount,
		ResetPassword,
	}
})
export default class Login extends Vue {
  private email: string;
  private password: string;
  private snackbar: boolean;

  constructor() {
    super();
    this.email = '';
    this.password = '';
    this.snackbar = false;
  }

  get ShowError(): boolean {
    const e: string = this.$store.getters.getLoginError;
    this.snackbar = e !== '';
    return e !== '';
  }

  set ShowError(show: boolean) {
    this.snackbar = show;
  }

  get LoginError(): string {
    return this.$store.getters.getLoginError;
  }

  get Account(): any {
    return this.$store.getters.getAccount;
  }

  computed() {
    return {
      email: '',
      password: '',
      snackbar: false,
    }
  }

  async handleSubmit() {
    await this.$store.dispatch('login', { 
      email: this.email, 
      password: this.password,
    });
    if (this.$store.getters.getAccount !== null) {
      this.$router.push('/account');
    }
  }
}
</script>

<style scoped lang='scss'>
.login {
  margin: 0 auto;
  flex: 1;

  .elevation-12 {
    flex: 1;

    .primary.v-btn:not(.v-btn--flat):not(.v-btn--text):not(.v-btn--outlined) {
      background-color: #4db6ac;
    }
  }

  .social {
    margin-top: 10px;
    & > div {
      padding: 10px;
      & > a {
        text-decoration: none;
        color: #2c3e50;

        & > button {
          display: flex;
          align-items: center;
          justify-content: center;
          border: 2px #eaeaea solid;
          border-radius: 4px;
          padding: 5px 10px;
          transition: 500ms all;
          
          & > span{
            padding: 0 5px 0 0;
          }
          
          & > i {
            padding: 0 0 5px 0;
          }
        }
        & > button:hover{
          box-shadow: 1px 0px 3px #000000;
        }
      }
    }
  }

  .create-account {
    padding: 1em 0;
    font-size: 1.3em;
    display: flex;
    flex-direction: row;
    justify-content: center;

    & > a, & > span {
      padding: 0.2em;
    }
    
    & > a {
      color: #4db6ac;
    }
  }
}
</style>
