package main

import (
  "math/rand"
  "time"
)

var (
  adjectives = [...]string{"Autumn","Hidden","Bitter","Misty","Silent","Empty","Dry","Dark","Summer","Icy","Delicate","Quiet","White","Cool","Spring","Winter","Patient","Twilight","Dawn","Crimson","Wispy","Weathered","Blue","Billowing","Broken","Cold","Damp","Falling","Frosty","Green","Long","Late","Lingering","Bold","Little","Morning","Muddy","Old","Red","Rough","Still","Small","Sparkling","Throbbing","Shy","Wandering","Withered","Wild","Black","Holy","Solitary","Fragrant","Aged","Snowy","Proud","Floral","Restless","Divine","Polished","Purple","Lively","Nameless","Puffy","Fluffy","Calm","Young","Golden","Avenging","Ancestral","Ancient","Argent","Reckless","Daunting","Short","Rising","Strong","Timber","Tumbling","Silver","Dusty","Celestial","Cosmic","Crescent","Double","Far","Half","Inner","Milky","Northern","Southern","Eastern","Western","Outer","Terrestrial","Huge","Deep","Epic","Titanic","Mighty","Powerful"}
  nouns = [...]string{"Waterfall","River","Breeze","Moon","Rain","Wind","Sea","Morning","Snow","Lake","Sunset","Pine","Shadow","Leaf","Dawn","Glitter","Forest","Hill","Cloud","Meadow","Glade","Bird","Brook","Butterfly","Bush","Dew","Dust","Field","Flower","Firefly","Feather","Grass","Haze","Mountain","Night","Pond","Darkness","Snowflake","Silence","Sound","Sky","Shape","Surf","Thunder","Violet","Wildflower","Wave","Water","Resonance","Sun","Wood","Dream","Cherry","Tree","Fog","Frost","Voice","Paper","Frog","Smoke","Star","Sierra","Castle","Fortress","Tiger","Day","Sequoia","Cedar","Wrath","Blessing","Spirit","Nova","Storm","Burst","Protector","Drake","Dragon","Knight","Fire","King","Jungle","Queen","Giant","Elemental","Throne","Game","Weed","Stone","Apogee","Bang","Cluster","Corona","Cosmos","Equinox","Horizon","Light","Nebula","Solstice","Spectrum","Universe","Magnitude","Parallax"}
)

func GenerateRandomName() string {
  n := nouns[rand.Intn(len(nouns))]
  a := adjectives[rand.Intn(len(adjectives))]
  return a + " " + n
}

// seeding the rand is global... yikes!
func init() {
  rand.Seed(time.Now().UTC().UnixNano())
}
