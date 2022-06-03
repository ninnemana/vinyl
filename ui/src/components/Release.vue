<template>
	<div class='root'>
		<div v-if='result' class="meta">
			<h1>{{ result.title }}</h1>
            <div class="photo">
                <Carousel :images='result.images' />
            </div>
        </div>
		<div v-else>No registered records. Maybe explore for something new?</div>
        <div>{{ result }}</div>
	</div>
</template>

<script lang='ts'>
import { Component, Prop, Vue } from 'vue-property-decorator';
import Carousel from '@/components/Carousel.vue';

@Component({
	components: {
		Carousel,
	}
})
export default class Release extends Vue {    
    get result() {
        return this.$store.getters.getResult || {};
    }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style scoped lang='scss'>
.feed {
	max-width: 1335px;
	margin: 0 auto;

	.results {
		display: flex;
		flex-flow: row wrap;
		justify-content: flex-start;

		.result {
			height: 350px;
			flex-basis: 20%;
			width: 259px;
			position: relative;
			padding: 10px;
			box-sizing: border-box;
			margin: 10px;
		}
	}
}
</style>
