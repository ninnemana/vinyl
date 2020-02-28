<template>
  <v-content>
    <v-container fluid>
      <v-row>
        <v-col class='login' cols='8'>
          <v-row :align='alignment' :justify='justify'>
            <v-card class='elevation-12'>
              <v-card-text>
                <v-form ref='form' v-model='valid' :lazy-validation='lazy'>
                  <v-text-field label='Login' name='login' prepend-icon='person' type='text' />

                  <v-text-field
                    id='password'
                    label='Password'
                    name='password'
                    prepend-icon='lock'
                    type='password'
                  />
                </v-form>
              </v-card-text>
              <v-card-actions>
                <v-spacer />
                <v-btn color='primary' class>Login</v-btn>
              </v-card-actions>
            </v-card>
          </v-row>
        </v-col>
      </v-row>
    </v-container>
  </v-content>
</template>

<script>
export default {
  name: 'Login',
  components: {},
  props: {
      valid: {},
  }
};
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
}
</style>
