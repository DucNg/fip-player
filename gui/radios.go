package gui

import (
	"encoding/base64"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
)

const baseURL = "https://www.radiofrance.fr/_app/remote/di23tz/getLive?payload="

func forgePayload(brandName string) string {
	json := fmt.Sprintf("[{\"brandName\":1},\"%v\"]", brandName)

	return base64.StdEncoding.EncodeToString([]byte(json))
}

var radios = []list.Item{
	item{
		title:       "France Culture",
		desc:        "S'informer, développer sa curiosité et apprendre chaque jour.",
		streamUrl:   "https://stream.radiofrance.fr/franceculture/franceculture.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("franceculture"),
	},
	item{
		title:       "France Inter",
		desc:        "Première radio d'actualité'généraliste et culturelle.",
		streamUrl:   "https://stream.radiofrance.fr/franceinter/franceinter.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("franceinter"),
	},
	item{
		title:       "franceinfo",
		desc:        "franceinfo, l'offre d'information en continu du service public.",
		streamUrl:   "https://stream.radiofrance.fr/franceinfo/franceinfo.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("franceinfo"),
	},
	item{
		title:       "FIP",
		desc:        "La radio la plus éclectique du monde",
		streamUrl:   "https://stream.radiofrance.fr/fip/fip.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("fip"),
	},
	item{
		title:       "FIP Cultes",
		desc:        "En illimité, une webradio dédiée uniquement aux titres cultes qui ont fait l’histoire et continuent de faire vibrer le présent.",
		streamUrl:   "https://stream.radiofrance.fr/fipcultes/fipcultes_hifi.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("fip_cultes"),
	},
	item{
		title:       "FIP Électro",
		desc:        "Un mix de titres inédits, de nouveautés et de grands classiques comme Cassius, Superpoze, Tosca.",
		streamUrl:   "https://stream.radiofrance.fr/fipelectro/fipelectro_hifi.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("fip_electro"),
	},
	item{
		title:       "FIP Groove",
		desc:        "Un mix R'n'B, Funk, Soul, Trip-Hop. Gardez le rythme avec Gil Scott Heron, De la Soul, Amy Winehouse, Marvin Gaye...",
		streamUrl:   "https://stream.radiofrance.fr/fipgroove/fipgroove_hifi.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("fip_groove"),
	},
	item{
		title:       "FIP Hip-Hop",
		desc:        "Découvrez une sélection éclectique de classiques et nouveautés du genre.",
		streamUrl:   "https://stream.radiofrance.fr/fiphiphop/fiphiphop_hifi.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("fip_hiphop"),
	},
	item{
		title:       "FIP Jazz",
		desc:        "Un mix de titres inédits et de grands classiques : d’Avishai Cohen à Herbie Hancock, de Nina Simone à Christian Scott.",
		streamUrl:   "https://stream.radiofrance.fr/fipjazz/fipjazz_hifi.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("fip_jazz"),
	},
	item{
		title:       "FIP Metal",
		desc:        "Une sélection éclectique à travers les genres : thrash metal, hard rock, death metal...",
		streamUrl:   "https://stream.radiofrance.fr/fipmetal/fipmetal_hifi.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("fip_metal"),
	},
	item{
		title:       "FIP Monde",
		desc:        "D’Ibrahim Maalouf à Tony Allen, de Vaudou Game à Bebel Gilberto, explorez le meilleur de la World Music.",
		streamUrl:   "https://stream.radiofrance.fr/fipworld/fipworld_hifi.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("fip_world"),
	},
	item{
		title:       "FIP Nouveautés",
		desc:        "Nouveaux artistes, dernières sorties d’album : ne manquez rien de l’actualité musicale.",
		streamUrl:   "https://stream.radiofrance.fr/fipnouveautes/fipnouveautes_hifi.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("fip_nouveautes"),
	},
	item{
		title:       "FIP Pop",
		desc:        "Écoutez une sélection variée d’artistes Pop.",
		streamUrl:   "https://stream.radiofrance.fr/fippop/fippop_hifi.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("fip_pop"),
	},
	item{
		title:       "FIP Reggae",
		desc:        "Écoutez le meilleur du reggae : roots, ska, dub, kumina.",
		streamUrl:   "https://stream.radiofrance.fr/fipreggae/fipreggae_hifi.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("fip_reggae"),
	},
	item{
		title:       "FIP Rock",
		desc:        "Un mix de nouveautés et de grands classiques comme David Bowie, Radiohead, Lou Reed...",
		streamUrl:   "https://stream.radiofrance.fr/fiprock/fiprock_hifi.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("fip_rock"),
	},
	item{
		title:       "FIP Sacré français !",
		desc:        "Découvrez le meilleur de la musique francophone gratuitement et sans pub : de Brassens à Orelsan, du Canada à l'Afrique, écoutez la sélection à la française des programmateurs de Fip.",
		streamUrl:   "https://stream.radiofrance.fr/fipsacrefrancais/fipsacrefrancais_hifi.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("fip_sacre_francais"),
	},
	item{
		title:       "France Musique",
		desc:        "Tout l'univers de la musique Classique, Jazz, Contemporaine et plus.",
		streamUrl:   "https://stream.radiofrance.fr/francemusique/francemusique_hifi.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("francemusique"),
	},
	item{
		title:       "France Musique - La Baroque",
		desc:        "Un voyage entre le XVIIe et le milieu du XVIIIe siècle.",
		streamUrl:   "https://stream.radiofrance.fr/francemusiquebaroque/francemusiquebaroque.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("francemusique_baroque"),
	},
	item{
		title:       "France Musique - Classique Easy",
		desc:        "La meilleure radio classique : les plus belles œuvres.",
		streamUrl:   "https://stream.radiofrance.fr/francemusiqueeasyclassique/francemusiqueeasyclassique_hifi.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("francemusique_classique_easy"),
	},
	item{
		title:       "France Musique - Classique Plus",
		desc:        "Des œuvres rares, intégrales et peu connues.",
		streamUrl:   "https://stream.radiofrance.fr/francemusiqueclassiqueplus/francemusiqueclassiqueplus.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("francemusique_classique_plus"),
	},
	item{
		title:       "France Musique - Concerts Radio France",
		desc:        "Des interprétations exceptionnelles de musique classique et symphonique.",
		streamUrl:   "https://stream.radiofrance.fr/francemusiqueconcertsradiofrance/francemusiqueconcertsradiofrance.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("francemusique_concert_rf"),
	},
	item{
		title:       "France Musique - La Contemporaine",
		desc:        "Découvrez des œuvres expérimentales et classiques, dans une sélection inédite réalisée par France Musique.",
		streamUrl:   "https://stream.radiofrance.fr/francemusiquelacontemporaine/francemusiquelacontemporaine_hifi.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("francemusique_la_contemporaine"),
	},
	item{
		title:       "France Musique - La Jazz",
		desc:        "Une sélection variée de titres jazz.",
		streamUrl:   "https://stream.radiofrance.fr/francemusiquelajazz/francemusiquelajazz_hifi.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("francemusique_la_jazz"),
	},
	item{
		title:       "France Musique - Musique de Films",
		desc:        "Écoutez les musiques de films les plus géniales avec la radio Musique de Films de France Musique, une sélection unique de musiques de cinéma connues et plus rares.",
		streamUrl:   "https://stream.radiofrance.fr/francemusiquelabo/francemusiquelabo_hifi.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("francemusique_evenementielle"),
	},
	item{
		title:       "France Musique - Ocora Musiques du monde",
		desc:        "Des pépites musicales traditionnelles du monde entier.",
		streamUrl:   "https://stream.radiofrance.fr/francemusiqueocoramonde/francemusiqueocoramonde.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("francemusique_ocora_monde"),
	},
	item{
		title:       "France Musique - Opéra",
		desc:        "Une sélection inédite des plus grands Opéras, des plus beaux airs, des plus belles voix.",
		streamUrl:   "https://stream.radiofrance.fr/francemusiqueopera/francemusiqueopera_hifi.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("francemusique_opera"),
	},
	item{
		title:       "France Musique - Piano Zen",
		desc:        "Vous voulez vous détendre ? Écoutez notre webradio « Piano Zen » qui explore plusieurs styles de piano, du classique au jazz. Alors, évadez-vous, laissez-vous porter et la musique… vous relaxer.",
		streamUrl:   "https://stream.radiofrance.fr/francemusiquepianozen/francemusiquepianozen_hifi.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("francemusique_piano_zen"),
	},
	item{
		title:       "Mouv'",
		desc:        "Tous les titres, hits et découvertes du rap, R'n'B, reggae & dancehall.",
		streamUrl:   "https://stream.radiofrance.fr/mouv/mouv.m3u8?id=radiofrance",
		metadataUrl: baseURL + forgePayload("mouv"),
	},
}
