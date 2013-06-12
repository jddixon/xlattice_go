// mt19937-64 with license removed for brevity

package rnglib

const NN int = 312
const MM int = 156
const MATRIX_A uint64 = 0xB5026F5AA96619E9
const UM uint64 = 0xFFFFFFFF80000000 /* Most significant 33 bits */
const LM uint64 = 0x7FFFFFFF         /* Least significant 31 bits */

type MT64 struct {
	mt  [NN]uint64 // state vector
	mti int
}

func NewMT64() *MT64 {
	return &MT64{mti: NN + 1} // means mt[NN] is not initialized
}

/* initializes mt[NN] with a seed */
func (m *MT64) init_genrand64(seed uint64) {
	m.mt[0] = seed
	for m.mti = 1; m.mti < NN; m.mti++ {
		m.mt[m.mti] = uint64(6364136223846793005)*(m.mt[m.mti-1]^(m.mt[m.mti-1]>>62)) + uint64(m.mti)
	}
}

/* initialize by an array with array-length */
/* init_key is the array for initializing keys */
/* key_length is its length */
func (m *MT64) init_by_array64(init_key []uint64, key_length uint64) {
	var i, j, k uint64
	const NN64 = uint64(NN)
	m.init_genrand64(uint64(19650218))
	i = 1
	j = 0
	if NN64 > key_length {
		k = NN64
	} else {
		k = key_length
	}
	for ; k > 0; k-- {
		m.mt[i] = (m.mt[i] ^ ((m.mt[i-1] ^ (m.mt[i-1] >> 62)) * 3935559000370003845)) + init_key[j] + j /* non linear */
		i++
		j++
		if i >= NN64 {
			m.mt[0] = m.mt[NN-1]
			i = 1
		}
		if j >= key_length {
			j = 0
		}
	}
	for k = NN64 - 1; k > 0; k-- {
		m.mt[i] = (m.mt[i] ^ ((m.mt[i-1] ^ (m.mt[i-1] >> 62)) * 2862933555777941757)) - i /* non linear */
		i++
		if i >= NN64 {
			m.mt[0] = m.mt[NN64-1]
			i = 1
		}
	}

	m.mt[0] = uint64(1) << 63 // MSB is 1; assuring non-zero initial array
}

/* generates a random number on [0, 2^64-1]-interval */
func (m *MT64) genrand64_int64() uint64 {
	var (
		i int
		x uint64
	)
	var mag01 = [2]uint64{0, MATRIX_A}
	// const NN64 = uint64(NN)

	if m.mti >= NN { /* generate NN words at one time */

		/* if init_genrand64() has not been called, */
		/* a default initial seed is used     */
		if m.mti == NN+1 {
			m.init_genrand64(uint64(5489))
		}

		for i = 0; i < NN-MM; i++ {
			x = (m.mt[i] & UM) | (m.mt[i+1] & LM)
			m.mt[i] = m.mt[i+MM] ^ (x >> 1) ^ mag01[int((x&uint64(1)))]
		}
		for ; i < NN-1; i++ {
			x = (m.mt[i] & UM) | (m.mt[i+1] & LM)
			m.mt[i] = m.mt[i+(MM-NN)] ^ (x >> 1) ^ mag01[(int)(x&uint64(1))]
		}
		x = (m.mt[NN-1] & UM) | (m.mt[0] & LM)
		m.mt[NN-1] = m.mt[MM-1] ^ (x >> 1) ^ mag01[(int)(x&uint64(1))]

		m.mti = 0
	}
	x = m.mt[m.mti]
	m.mti++

	x ^= (x >> 29) & uint64(0x5555555555555555)
	x ^= (x << 17) & uint64(0x71D67FFFEDA60000)
	x ^= (x << 37) & uint64(0xFFF7EEE000000000)
	x ^= (x >> 43)

	return x
}

/* generates a random number on [0, 2^63-1]-interval */
func (m *MT64) genrand64_int63() uint64 {
	return uint64(m.genrand64_int64() >> 1)
}

/* generates a random number on [0,1]-real-interval */
func (m *MT64) genrand64_real1() float64 {
	return float64(m.genrand64_int64()>>11) * (1.0 / 9007199254740991.0)
}

/* generates a random number on [0,1)-real-interval */
func (m *MT64) genrand64_real2() float64 {
	return float64(m.genrand64_int64()>>11) * (1.0 / 9007199254740992.0)
}

/* generates a random number on (0,1)-real-interval */
func (m *MT64) genrand64_real3() float64 {
	return (float64(m.genrand64_int64()>>12) + 0.5) * (1.0 / 4503599627370496.0)
}

// rand.Source interface ////////////////////////////////////////////

func (m *MT64) Seed(seed int64) {
	m.init_genrand64(uint64(seed))
}
func (m *MT64) Int63() int64 {
	return int64(m.genrand64_int63())
}
