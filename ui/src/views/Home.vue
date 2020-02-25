<template>
	<div class='home'>
		<img alt='VinylTap Logo' src='../assets/logo.png' />
		<Searchbox />
	</div>
</template>

<script>
// @ is an alias to /src
import Searchbox from '@/components/Searchbox.vue';

export default {
	name: 'Home',
	components: {
		Searchbox
	}
};
</script>
