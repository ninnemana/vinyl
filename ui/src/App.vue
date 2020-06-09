<template>
  <v-app>
    <v-app-bar app short dark>
      <v-app-bar-nav-icon></v-app-bar-nav-icon>
      <v-toolbar-title>
        <router-link to='/'>Vinylshare.io</router-link>
      </v-toolbar-title>
      <v-spacer></v-spacer>
      <v-btn icon>
        <v-icon>mdi-magnify</v-icon>
      </v-btn>

      <v-menu left bottom>
        <template v-slot:activator="{ on, attrs }">
          <v-btn
            text
            v-bind="attrs"
            v-on="on"
          >
            <v-icon>mdi-dots-vertical</v-icon>
          </v-btn>
        </template>

        <AuthorizationList :authorized='Authed' />
      </v-menu>
    </v-app-bar>
    <v-sheet class="body">
      <router-view/>
    </v-sheet>
  </v-app>
</template>

<script lang='ts'>
// @ is an alias to /src
import AuthorizationList from '@/components/AuthorizationList.vue'
import { Component, Vue } from 'vue-property-decorator';

@Component({
  components: {
    AuthorizationList,
  }
})
export default class App extends Vue {
    get Authed(): boolean {
      const acct = this.$store.getters.getAccount;
      if (acct === null || acct.user === undefined) {
        return false;
      }

      return true;
    }
}
</script>

<style lang='scss'>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

#nav {
  padding: 30px;

  a {
    font-weight: bold;
    color: #2c3e50;

    &.router-link-exact-active {
      color: #42b983;
    }
  }
}

.body {
  padding-top: 80px;
}
</style>
