// chris 052515 Common code for test routines.

package steg

import (
	"os"
	"strings"
	"testing"
	"time"

	"math/rand"
)

const helloString = "hello, there, how are you? fine."

// The byte embedded in the above string.
const helloByte = 0xdb

func testHelloChunk(atomSize uint8) *chunk {
	ctx := NewCtx(atomSize)
	c := ctx.newChunk()
	n := len(c.data) / len(helloString)
	if n*len(helloString) != len(c.data) {
		panic("non-integral chunk multiple")
	}
	copy(c.data, []byte(strings.Repeat(helloString, n)))
	return c
}

func testLoremChunk() *chunk {
	ctx := NewCtx(2)
	c := ctx.newChunk()
	if len(loremString) != 8192 {
		panic("short lorem")
	}
	copy(c.data, []byte(loremString))
	return c
}

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UTC().UnixNano())
	os.Exit(m.Run())
}

// The bytes embedded in the below string.
var loremBytes = []byte{0x2a, 0x41}

// 8KiB of Lorem Ipsum.  Ideal for atom size 2 tests.
const loremString = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas id rutrum arcu, et feugiat tellus. Sed mattis, ipsum euismod volutpat interdum, elit erat placerat erat, id tincidunt augue ex id eros. Sed rutrum velit nec sollicitudin faucibus. Sed id cursus tortor. Cras non malesuada ante. Ut sapien lectus, accumsan bibendum dictum non, lobortis ut sem. Aenean quis egestas quam. Donec bibendum rhoncus eros, quis consectetur ex iaculis id.
In hac habitasse platea dictumst. Donec sem mauris, pretium a pharetra in, pretium eget nisl. Mauris sem nunc, dapibus ac maximus vitae, volutpat ut orci. Etiam blandit porta lorem, in imperdiet justo aliquam nec. Mauris pellentesque ac lacus in scelerisque. Phasellus rutrum neque eget sagittis euismod. Donec sit amet metus vehicula nulla lacinia placerat eget sed augue. Vestibulum lacinia sed nisi vitae placerat. Ut vel nisi ac est egestas mollis id scelerisque elit. Curabitur at aliquam felis. Morbi semper maximus molestie. Vestibulum convallis ex a sapien consequat, nec maximus libero pellentesque. Suspendisse ornare efficitur pulvinar. Donec sed arcu a justo imperdiet pretium. Mauris a molestie tortor, a sollicitudin est.
Donec vel est placerat nisi sollicitudin placerat. In posuere in odio eu hendrerit. Fusce dictum fringilla suscipit. Nunc aliquet, urna nec pharetra gravida, velit metus dictum lectus, non pulvinar libero odio vel diam. Cras sodales laoreet mauris id tempus. In rhoncus eu lorem a iaculis. Donec purus sapien, lacinia sit amet justo ut, rhoncus scelerisque elit. Curabitur mattis tortor tellus, sit amet imperdiet libero ornare et. Vivamus pulvinar massa et nibh sodales, et malesuada augue consequat. Aliquam sodales tincidunt est, tincidunt aliquam diam maximus ut. Sed tincidunt sem sem, non finibus quam sollicitudin nec. Sed lorem enim, condimentum at accumsan eu, euismod eget nunc. Donec cursus nibh ac velit consectetur, in auctor nisl convallis.
Fusce blandit ante sapien. Sed in tortor metus. Ut suscipit tristique ligula nec fermentum. Curabitur ac rhoncus lacus. Proin nisl ipsum, dignissim vitae pulvinar non, posuere et nunc. Curabitur fermentum ac lacus nec fermentum. Nunc ultricies placerat mauris, at tempor tellus tempor quis. Praesent ultrices urna sed ultrices placerat. Nunc in vestibulum ligula.
Maecenas eget lorem sed est lacinia porta id ut mauris. Maecenas maximus libero mattis consectetur tincidunt. Phasellus malesuada eget nibh at fringilla. In semper massa felis, nec efficitur augue iaculis id. Pellentesque ut dignissim massa. Fusce a sapien hendrerit, volutpat nulla sed, ornare nunc. Sed et sollicitudin risus. Mauris at metus ac risus varius sollicitudin. Fusce tempor tempus lorem et hendrerit.
Aliquam auctor quis dui non interdum. Sed elementum magna non tempor iaculis. Pellentesque at aliquam urna. Nulla sagittis, nulla a dictum vulputate, mi elit consequat odio, at tempus neque urna et ex. Aliquam ornare ante sit amet porttitor molestie. Cras a augue laoreet tellus vulputate posuere. Aliquam erat volutpat. Nulla suscipit feugiat consequat. Morbi porttitor leo vitae diam euismod euismod. Nulla facilisi. Mauris iaculis condimentum sodales. Maecenas et tellus et nulla pulvinar finibus. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Proin magna nisl, dapibus vel malesuada sed, semper sed neque. Integer fringilla est in urna placerat varius. Nulla vulputate ligula non massa ullamcorper, sed accumsan orci placerat.
Mauris mattis ipsum nec lobortis vestibulum. Praesent urna lorem, sodales quis tincidunt sit amet, eleifend congue sapien. Vivamus volutpat mi a lectus finibus molestie. Quisque pulvinar cursus volutpat. Donec condimentum risus justo, vel fermentum diam accumsan at. Vestibulum aliquam metus ut vestibulum mattis. Morbi laoreet orci nec orci consectetur lacinia. Quisque mollis enim eu diam sagittis, sit amet convallis lectus imperdiet. Ut velit libero, efficitur eu neque vel, lacinia fermentum eros. Vivamus consectetur nisl id lacus lobortis facilisis id a quam. Nullam auctor pretium neque, at varius odio venenatis a. Aliquam aliquet mi vitae hendrerit iaculis. Aenean interdum tellus et risus vehicula bibendum. Mauris id gravida est. Maecenas varius, augue sit amet dapibus tincidunt, ante felis dictum lorem, quis interdum nunc lacus vel ex. Sed a sapien dictum, cursus arcu eu, vestibulum velit.
In at egestas leo. Integer mattis eros in neque consectetur convallis. Sed sit amet turpis eu sapien vestibulum molestie. Quisque vel nulla at enim euismod varius eget vel nibh. Nam egestas, nibh in luctus feugiat, ligula dui vulputate odio, vestibulum fringilla sapien lacus id ipsum. Pellentesque ac ultrices erat. Aliquam luctus eget justo nec tempor. Vivamus porttitor sapien in lectus blandit, at iaculis metus dignissim. Proin fermentum facilisis tortor, sed tincidunt elit cursus ut.
Suspendisse sit amet elit nec sapien feugiat hendrerit. Ut erat ligula, tincidunt id ex eu, scelerisque consequat mi. Aenean viverra augue eget faucibus rhoncus. Nulla a neque vehicula, placerat ligula eu, sollicitudin ligula. Sed eros nisl, sollicitudin quis enim a, luctus lacinia tellus. Donec in mi diam. Nulla sit amet nisl ut mi placerat volutpat. Ut faucibus mi lacus, eu mollis orci venenatis et. Donec sollicitudin metus sit amet laoreet pellentesque. Aliquam eget diam blandit, vehicula sem auctor, consequat metus. Quisque imperdiet leo in vestibulum tincidunt. Integer consectetur ex eget diam sollicitudin, ut aliquet justo lobortis.
Sed accumsan elit vitae odio rutrum scelerisque. Mauris a iaculis velit. Cras posuere posuere risus, in vulputate justo euismod ac. Sed scelerisque orci nec tortor feugiat, eu interdum neque finibus. Sed commodo elit enim. Praesent sit amet sem sagittis arcu ornare pharetra eget id ipsum. In id dolor lacinia, volutpat ante id, auctor eros. Mauris sed ante ac urna bibendum rhoncus a sed mauris. Ut ligula sem, semper eu tincidunt non, rhoncus vel dui.
Nunc ultrices est eu nunc ultricies, id vestibulum libero sodales. Praesent congue porta odio, et bibendum ante cursus nec. Ut turpis justo, tristique vehicula ante eget, malesuada efficitur est. Morbi finibus ex ut dapibus interdum. Nam hendrerit magna vel dignissim ornare. In vehicula consequat nulla. Pellentesque at mauris laoreet, semper leo at, congue lorem. Mauris volutpat, lectus vel convallis interdum, sem nibh scelerisque tellus, eu euismod metus ante in felis. Etiam malesuada semper lectus. Mauris nec erat dignissim, tristique risus congue, tincidunt justo. Quisque euismod semper semper. Maecenas non risus eget tellus cursus hendrerit vel at ex. Sed at lacus ante. Praesent ac urna porta nulla dictum ultricies. Pellentesque blandit mi ipsum, in aliquet tortor viverra nec.
Donec sit amet urna lorem. Vivamus felis orci, feugiat id sodales ut, sodales ut velit. Nulla rhoncus dictum vehicula. Quisque egestas sapien ligula, et faucibus nisl dictum in. Maecenas gravida ornare elit sit amet tempor. Nam malesuada diam lacus, sed faucibus eros mollis ac. Donec aliquam diam et est eleifend congue. Vestibulum ut elementum eros. Integer feugiat placerat ipsum nec consequat. Vivamus ut ligula neque. Proin quis auctor odio. Sed ultrices dictum massa, non vulputate enim euismod id. Aenean vitae egestas leo. Donec tincidunt congue lorem, vel tempus augue volutpat ut. Sed nec rutrum ligula, congue semper odio.
Aenean eget rhoncus erat, ut tristique nisl. Nunc non pellentesque neque. Mauris metus odio, euismod sit amet eros nec, ultrices tempor augue. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Maecenas ac nunc et dolor venenatis convallis. Vestibulum sed faucibus elit. Cras euismod at augue sed varius.
Suspendisse potenti. Aliquam placerat dolor facilisis est sodales convallis. Etiam aliquet tempus erat quis facilisis. Vestibulum ut libero nec risus ullamcorper posuere. Nunc sit amet tempus nisl, id lacinia dui. Donec sit amet nisl eget erat fermentum accumsan nec at est. Etiam rhoncus sed enim sit amet malesuada.
Mauris nec aliquam erat, sed tincidunt nulla. Pellentesque id suscipit dolor. Sed a eros tristique, rhoncus est eget.`
