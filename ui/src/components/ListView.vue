<template>
	<div class="entry">
		<p>Kendrick Lamar</p>
        <p>Good Kid Madd City</p>
	</div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component
export default class ListView extends Vue {
	@Prop() private entry = {};
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.entry{
    display: flex;
    flex-direction: row;
    justify-content: space-around;
}
</style>
