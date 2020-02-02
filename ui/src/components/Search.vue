<template>
	<form class="search" @submit.prevent="handleSubmit">
		<input type="search" v-model="query" placeholder="Discover something new .." />
	</form>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component
export default class Search extends Vue {
    private query = '';

    beforeEnter(to: any) {
        this.$store.dispatch('search', to.params.query)
    }
    
    private handleSubmit() {
        this.$store.dispatch('search', this.query);
        this.$router.push(`/search?query=${this.query}`);
    }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.search{
    
    input[type=search]{
        border: none;
        background: none;
        border-bottom: .4rem solid gray;
        padding: .5rem .2rem;
        margin: 1rem 0;
        width: 100%;
        font-size: 2rem;

        &:focus{
            outline: none;
        }
    }
}
</style>
