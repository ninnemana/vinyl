<template>
	<div v-bind:key='entry.id' class='entry'>
		<img v-bind:src='entry.thumb' v-bind:alt='entry.title' />
		<span>{{entry.title}}</span>
		<router-link :to='{ path: `/release/${encodeURIComponent(entry.title)}/${entry.id}`}'>details</router-link>
	</div>
</template>

<script lang='ts'>
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component
export default class ListView extends Vue {
	@Prop() private entry!: {};
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style scoped lang='scss'>
.entry {
	display: flex;
	flex-direction: column;
	justify-content: center;
	padding: 5px;
	-webkit-box-sizing: initial;
	-moz-box-sizing: initial;
	box-sizing: initial;
	background: #ececec;
	margin: 0;
	height: 100%;
	width: 100%;
	overflow: hidden;
	-webkit-transition: padding 0.15s cubic-bezier(0.4,0,0.2,1), margin 0.15s cubic-bezier(0.4,0,0.2,1), box-shadow 0.15s cubic-bezier(0.4,0,0.2,1);
	transition: padding 0.15s cubic-bezier(0.4,0,0.2,1), margin 0.15s cubic-bezier(0.4,0,0.2,1), box-shadow 0.15s cubic-bezier(0.4,0,0.2,1);
	position: relative;
}

.column {
	flex-basis: 100%;
	margin: 10px;
}
</style>
