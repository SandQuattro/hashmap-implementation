package main

import (
	"encoding/binary"
	"fmt"
)

// требования к map
// должна быть hash функция bucket = hash(key)
// hash функция - одностороннее преобразование, позволяющее получить индекс бакета для получения / вставки значения
// hash функция должна обладать следующими свойствами
// 1. Равномерность. все записи должны быть равномерно распределены по бакетам
// 2. Быстрота. Доступ должен стремиться к константному (O)1
// 3. Детерминированность. Для одного и того же ключа должно быть одно и то же значение бакета. Положили в бакет по ключу
// и должны получить по этому ключу это же значение
// 4. безопасность. Чтобы нельзя было подобрать ключи так, что они все попали в один бакет и тогда скорость доступа станет (O)n

// Определим структуру заголовка hashmap
type hmap struct {
	// количество элементов в map
	size uint8

	// почему бакеты хранятся в виде логарифма?
	// 1) логарифм позволяет хранить более маленькое число, что позволяет немного сэкономить на памяти
	// 2) ускоряет побитовые операции, которые проводим с данным полем
	buckets uint8 // log_2, math.Log2(n)

	// *Buckets указатель на список бакетов, где каждому бакету соответствуют младшие биты хеша (LowOrderBits). Их назначение
	// помогать искать бакеты из предоставленной хеш функции

	// Выполнение условия безопасности, предъявляемое к hashmap, чтобы нельзя было подобрать ключи так,
	// чтобы они все попали в один бакет и тогда скорость доступа станет (O)n
	// seed uint32
}

var h = 2232323424 // 5461

func hashfunc(key string) uint {
	// здесь мы должны выполнить некое преобразование ключа, получив на выходе число, которое будет хешем

	// для примера
	hash := h
	return uint(hash)
}

// Принимает один аргумент b типа uint8 и возвращает значение типа uintptr.
// Функция используется для получения маски, вычисляющей младшие биты хеша, которая требуется для
// нахождения индекса в хэш-таблице.
// Cначала вызывается bucketShift(b), чтобы получить 2^b(количество бакетов), и затем из результата вычитается 1.
// Это делается для создания битовой маски (например, если b = 3, то результат будет 111 в двоичном виде,
// что эквивалентно 7 в десятичном).
// пример: 000000100 - 1 = 1 пропадает в 3м рязряде и предыдущие разряды заполняются 1: 000000011
func bucketMask(shift uintptr) uintptr {
	return shift - 1
}

// bucketShift принимает аргумент b типа uint8 - количество бакетов, которое хранится как log_2 и возвращает значение
// типа uintptr.
// Смысл этой функции в том, чтобы возвратить 2 в степени b.
// пример: у нас 4 бакета, они хранятся как log^2 = log(4) = 2
// берем 1 = 000000001 сдвигаем ее на log(4) = 2
// получим 000000100
func bucketShift(b uint8) uintptr {
	return 1 << b
}

func main() {
	// для примера создадим map с 4 бакетами
	m := hmap{
		size:    0,
		buckets: 2, // log_2. n бакетов, 2^n, 2^2 = 4
	}
	_ = m

	// получаем hash
	hash := hashfunc("key")

	// затем по хешу мы должны получить номер бакета
	// в общем случае бакет получается как остаток от деления хеша на количество бакетов

	// для примера попробуем получить номер бакета этим способом
	bucket := hash % 4
	fmt.Printf("got bucket number using mod method: %d\n", bucket)
	// Получаем номер бакета = 1

	// Но для ускорения вычислений операция получения остатка от деления выполняется побитово
	// 1) Приводим хеш к байтовому представлению
	fmt.Printf("hash: %d, binary representation of hash: %b\n", hash, hash)

	// 2) Нам потребуется логарифм от количества бакетов, который хранится в заголовке map
	fmt.Printf("buckets (as log_2): %d\n", m.buckets)

	// Например, количество бакетов = 4, тогда log_2 = 2, передаем b = 2 в функцию bucketShift
	// Выполняем побитовый сдвиг 1 << b
	// 1 в двоичном виде это 0001, убедимся в этом
	i := make([]byte, 8)
	i[0] = 1
	fmt.Printf("binary representation of 1: %b %b\n", binary.BigEndian.Uint16(i), binary.LittleEndian.Uint16(i))

	// Понимание Little Endian и Big Endian
	// Little Endian: В этом формате младший байт (байт с наименьшим значащим битом) хранится в адресе с наименьшим числом (то есть на "начале" числа).
	// Большие эндианы хранят данные точно наоборот.
	// Big Endian: В этом формате старший байт (байт с наибольшим значащим битом) хранится в адресе с наименьшим числом (на "начале" числа).
	// Влияние на bucketShift и bucketMask
	// Операции, выполняемые функциями bucketShift и bucketMask, осуществляют сдвиги и маскирование на уровне битов,
	// не обращая внимания на то, как эти биты могут быть упакованы в байты и как эти байты выравнены в памяти.
	// В результате, порядок endian здесь не играет роли.
	//
	// Last Order Bits: Когда мы говорим о "last order bits" или последних значащих битах из контекста этих функций, мы просто обращаем внимание на биты, которые определяются после операции сдвига. Эти биты определяют значение без учета того, как они представлены в памяти в виде байтов.
	// Короче говоря, когда вы работаете с этими операциями на битовом уровне для таких задач, как вычисление бакетов в хэш-таблице или маскирование значений, порядок байтов Little Endian или Big Endian не изменяет логику вычислений. Ваши операции сдвига и маскирования будут работать одинаково вне зависимости от порядка байтов.
	// выполняем сдвиг единички влево на количество бакетов в виде логарифма log_2 = 2 - получается на 2 позиции влево
	// было - 0000001, стало - 0000100 - это 4 в десятичной системе
	shift := bucketShift(m.buckets)
	fmt.Printf("shift value result:%d\n", shift)

	// и теперь от количества бакетов 4(0000100) отнимаем 1 = 3(0000011)
	mask := bucketMask(shift)
	fmt.Printf("mask result:%d, binary representation: %b\n\n", mask, mask)

	// В результате всех манипуляций мы получили битовую маску, для получения кусочка двоичного числа, в данном случае получение LOB
	// В нашем случае значение хеша = 1010101010101 и мы хотим получить LOB, это 2 последних бита
	// Для того, чтобы это сделать, воспользуемся битовой маской, полученной выше, у которой все разряды равны 0,
	// кроме тех, которые мы должны получить 0011.
	// Выполняем операцию AND, логическое ИЛИ
	// 1010101010101 AND 00000000011, AND вернет нам 1 в разряде, где есть 1, остальные забьет нулями
	// Наложив маску на хеш, получим LOB(hash) = 01 - это будет бакет, в который мы должны вставить значение
	result := hash & uint(bucketMask(shift))
	fmt.Printf("result: %d, binary representation: %b\n", result, result)

	// Сравним результат
	result = hashfunc("key") % 4
	fmt.Printf("result: %d, binary representation: %b\n", result, result)

	// Если бы количество бакетов было 8, то нам понадобилось бы 3 бита младшего порядка, log_3 = 2^3 = 8 и тд
	// LOB(hash) = 011 и тд

	// Там есть еще нюанс, мы должны знать размер бакета, чтобы корректно спозиционироваться внутри бакета
	// есть наш hash - значение хеша, и элементы в бакете 0[........]1[........]2[........]
	// таким образом нам надо перепрыгнуть через элементы и спозиционироваться
	// берем указатель на бакеты - m.buckets, (hash & mask)*uintptr(bucketsize) - смещение в бакете
	// По сути преобразовываем указатель в uintptr - по сути в число, прибавляем к нему наше смещение, определенное выше
	// (hash & mask)*uintptr(bucketsize) и возвращаем снова указатель на нужную нам позицию
}
