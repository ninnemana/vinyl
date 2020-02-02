<template>
  <div class="search">
    <h1>These are the search results.</h1>
    <div class="row">
      <div v-for="result in searchResults()" v-bind:key="result.id" class="column">
        <img v-bind:src=result.thumb v-bind:alt=result.title>
        <span>{{result.title}}</span>
        <router-link :to="{ path: `/search/${result.id}`}">details</router-link>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
// import Feed from '@/components/Feed.vue';
// import Search from '@/components/Search.vue';

@Component({
	components: {
		// Feed,
		// Search,
	},
})
export default class Search extends Vue {
  searchResults(): any {
    return this.$store.state.search_results;
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.search {
  width: 80%;
  max-width: 1200px;
  margin: 0 auto;
}

.row {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
}

.column {
  flex-basis: 100%;
  margin: 10px;
}

@media screen and (min-width: 800px) {
  .column {
    flex: 1;
  }
}
</style>

