package gui

import "github.com/charmbracelet/bubbles/list"

var radios = []list.Item{
	item{
		title:       "FIP",
		desc:        "La radio la plus éclectique du monde",
		streamUrl:   "https://stream.radiofrance.fr/fip/fip.m3u8?id=radiofrance",
		metadataUrl: "https://www.radiofrance.fr/api/v2.0/stations/fip/webradios/fip",
	},
	item{
		title:       "FIP Jazz",
		desc:        "Un mix de titres inédits et de grands classiques : d’Avishai Cohen à Herbie Hancock, de Nina Simone à Christian Scott.",
		streamUrl:   "https://stream.radiofrance.fr/fipjazz/fipjazz_hifi.m3u8?id=radiofrance",
		metadataUrl: "https://www.radiofrance.fr/api/v2.0/stations/fip/webradios/fip_jazz",
	},
	item{
		title:       "La contemporaine",
		desc:        "Découvrez des œuvres expérimentales et classiques, dans une sélection inédite réalisée par France Musique.",
		streamUrl:   "https://stream.radiofrance.fr/francemusiquelacontemporaine/francemusiquelacontemporaine_hifi.m3u8?id=radiofrance",
		metadataUrl: "https://www.radiofrance.fr/api/v2.0/stations/francemusique/webradios/francemusique_la_contemporaine",
	},
	item{
		title:       "FIP Groove",
		desc:        "Un mix R'n'B, Funk, Soul, Trip-Hop. Gardez le rythme avec Gil Scott Heron, De la Soul, Amy Winehouse, Marvin Gaye...",
		streamUrl:   "https://stream.radiofrance.fr/fipgroove/fipgroove_hifi.m3u8?id=radiofrance",
		metadataUrl: "https://www.radiofrance.fr/api/v2.0/stations/fip/webradios/fip_groove",
	},
	item{
		title:       "FIP Rock",
		desc:        "Un mix de nouveautés et de grands classiques comme David Bowie, Radiohead, Lou Reed...",
		streamUrl:   "https://stream.radiofrance.fr/fiprock/fiprock_hifi.m3u8?id=radiofrance",
		metadataUrl: "https://www.radiofrance.fr/api/v2.0/stations/fip/webradios/fip_rock",
	},
	item{
		title:       "FIP Metal",
		desc:        "Une sélection éclectique à travers les genres : thrash metal, hard rock, death metal...",
		streamUrl:   "https://stream.radiofrance.fr/fipmetal/fipmetal_hifi.m3u8?id=radiofrance",
		metadataUrl: "https://www.radiofrance.fr/api/v2.0/stations/fip/webradios/fip_metal",
	},
	item{
		title:       "FIP Hip-Hop",
		desc:        "Découvrez une sélection éclectique de classiques et nouveautés du genre.",
		streamUrl:   "https://stream.radiofrance.fr/fiphiphop/fiphiphop_hifi.m3u8?id=radiofrance",
		metadataUrl: "https://www.radiofrance.fr/api/v2.0/stations/fip/webradios/fip_hiphop",
	},
	item{
		title:       "FIP Pop",
		desc:        "Écoutez une sélection variée d’artistes Pop.",
		streamUrl:   "https://stream.radiofrance.fr/fippop/fippop_hifi.m3u8?id=radiofrance",
		metadataUrl: "https://www.radiofrance.fr/api/v2.0/stations/fip/webradios/fip_pop",
	},
	item{
		title:       "FIP Electro",
		desc:        "Un mix de titres inédits, de nouveautés et de grands classiques comme Cassius, Superpoze, Tosca.",
		streamUrl:   "https://stream.radiofrance.fr/fipelectro/fipelectro_hifi.m3u8?id=radiofrance",
		metadataUrl: "https://www.radiofrance.fr/api/v2.0/stations/fip/webradios/fip_electro",
	},
	item{
		title:       "FIP Monde",
		desc:        "D’Ibrahim Maalouf à Tony Allen, de Vaudou Game à Bebel Gilberto, explorez le meilleur de la World Music.",
		streamUrl:   "https://stream.radiofrance.fr/fipworld/fipworld_hifi.m3u8?id=radiofrance",
		metadataUrl: "https://www.radiofrance.fr/api/v2.0/stations/fip/webradios/fip_world",
	},
	item{
		title:       "FIP Reggae",
		desc:        "Écoutez le meilleur du reggae : roots, ska, dub, kumina.",
		streamUrl:   "https://stream.radiofrance.fr/fipreggae/fipreggae_hifi.m3u8?id=radiofrance",
		metadataUrl: "https://www.radiofrance.fr/api/v2.0/stations/fip/webradios/fip_reggae",
	},
	item{
		title:       "FIP Nouveautés",
		desc:        "Nouveaux artistes, dernières sorties d’album : ne manquez rien de l’actualité musicale.",
		streamUrl:   "https://stream.radiofrance.fr/fipnouveautes/fipnouveautes_hifi.m3u8?id=radiofrance",
		metadataUrl: "https://www.radiofrance.fr/api/v2.0/stations/fip/webradios/fip_nouveautes",
	},
}
