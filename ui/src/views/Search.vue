<template>
  <div class="search">
    <h1>These are the search results.</h1>
    <div v-for="result in searchResults()" v-bind:key="result">
      {{result}}
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
    return Array();
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.search {
	width: 80%;
	margin: auto;
}
</style>

