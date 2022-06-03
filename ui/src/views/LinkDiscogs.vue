<template>
    <v-content>
        <v-container fluid>
            <h1>Add Discogs Account</h1>
            <v-form>
                <v-container>
                    <v-row justify="center">
                        <v-col cols="12" md="8">
                            <v-text-field 
                                :counter="10" 
                                label="Name" 
                                required
                                center
                            />
                        </v-col>
                    </v-row>
                    <v-row justify="center">
                        <v-col cols="12" md="8">
                            <v-text-field 
                                label="Email address" 
                                required 
                            />
                        </v-col>
                    </v-row>
                </v-container>
            </v-form>
        </v-container>
    </v-content>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';

@Component({
	components: {}
})
export default class LinkDiscogs extends Vue {
    private account: any;

    constructor() {
        super();
        this.account = null;
    }

    computed() {
        return {
            account: this.$store.getters.getAccount,
        };
    }
}
</script>
