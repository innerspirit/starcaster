package main

// Do not change this name, nux need manifest to generate AndroidManifest.xml
const manifest = `
{
    import: {
        ui: "nuxui.org/nuxui/ui",
    },

    application: {
        // display name at luancher 
		label: starcaster,  

        // application identifier name
        name: "org.nuxui.samples.starcaster",
    },

    permissions: [
        // wifi,
        storage,
        viewPhoto,
        savePhoto,
    ],

    mainWindow: {
        width: 15%,
        height: 1:1,
        title: "StarCaster",
        content: {
            type: ui.Layer,
            width: 100%,
            height: 50%,
            children: [
                {
                    widget: ui.Image,
                    src: "starcaster.png",
                    type: ui.Image,
                    width: 100%,
                    height: 1:1,
                    margin: {top: 4wt, bottom: 3wt}
                },

            ]
        }
        background: #000000,
    },
}
`
