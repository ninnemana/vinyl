<template>
  <v-content>
    <v-container fluid>
      <v-row>
        <v-col class='login' cols='8'>
          <v-row align='center' justify='center'>
            <v-card class='elevation-12'>
              <v-card-text>
                <v-form @submit.prevent='handleSubmit' id='login-form' ref='form' lazy-validation>
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
                </v-form>
              </v-card-text>
              <v-card-actions>
                <v-spacer />
                <v-btn form='login-form' @click='handleSubmit' color='primary' class>Login</v-btn>
              </v-card-actions>
            </v-card>
          </v-row>
          <v-row align='center' justify='center' class='create-account'>
            <router-link to='/account/create'>Create Account</router-link>
            <span>|</span>
            <router-link to='/account/reset-password'>Reset Password</router-link>
          </v-row>
        </v-col>
      </v-row>
    </v-container>
  </v-content>
</template>

<script>
import { Component, Vue } from 'vue-property-decorator';

@Component({})
export default class ResetPassword extends Vue {
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
