package tempconv
import "fmt"
type Celsius float64
type Fahrenheit float64
type Kelvin float64
const (
AbsoluteZeroC Celsius = -273.15
FreezingC Celsius = 0
BoilingC Celsius = 100
AbsoluteZeroK Kelvin = -273.15
)

func CToF(c Celsius) Fahrenheit { return Fahrenheit(c*9/5 + 32) }
func FToC(f Fahrenheit) Celsius { return Celsius((f-32)* 5/9) }
func KtoC(k Kelvin) Celsius { return Celsius(k+723.15)}

func (c Celsius) String() string { return fmt.Sprintf("%g°C", c) }
func (f Fahrenheit) String() string { return fmt .Sprintf ("%g°F",f) }