/*
 * Copyright (c) 2008-2021, Hazelcast, Inc. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License")
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package hazelcast_test

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	hz "github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/internal/it"
	"github.com/hazelcast/hazelcast-go-client/types"
)

func TestMultiMap_Put(t *testing.T) {
	it.MultiMapTester(t, func(t *testing.T, m *hz.MultiMap) {
		ctx := context.Background()
		targetValue := "value1"
		success, err := m.Put(ctx, "key", targetValue)
		if err != nil {
			t.Fatal(err)
		}
		assert.True(t, success, "multi-map put failed")
		values, err := m.Get(ctx, "key")
		if err != nil {
			t.Fatal(err)
		}
		it.AssertEquals(t, []interface{}{targetValue}, values)
		targetValue2 := "value2"
		success, err = m.Put(ctx, "key", targetValue2)
		if err != nil {
			t.Fatal(err)
		}
		assert.True(t, success, "multi-map put failed")
		values, err = m.Get(ctx, "key")
		if err != nil {
			t.Fatal(err)
		}
		assert.ElementsMatch(t, []interface{}{targetValue, targetValue2}, values)
	})
}

func TestMultiMap_PutAll(t *testing.T) {
	it.MultiMapTester(t, func(t *testing.T, m *hz.MultiMap) {
		ctx := context.Background()
		existingValue := "v4"
		assert.True(t, it.MustBool(m.Put(ctx, "k1", existingValue)), "multimap put failed")
		value := it.MustValue(m.Get(ctx, "k1"))
		it.AssertEquals(t, []interface{}{existingValue}, value)
		if err := m.PutAll(ctx, "k1", "v1", "v2", "v3"); err != nil {
			t.Fatal(err)
		}
		assert.ElementsMatch(t, []string{"v1", "v2", "v3", "v4"}, it.MustSlice(m.Get(ctx, "k1")))
	})
}

func TestMultiMap_Delete(t *testing.T) {
	it.MultiMapTester(t, func(t *testing.T, m *hz.MultiMap) {
		ctx := context.Background()
		targetValue := "value"
		assert.True(t, it.MustBool(m.Put(ctx, "key", targetValue)), "multimap put failed")
		value := it.MustValue(m.Get(ctx, "key"))
		it.AssertEquals(t, []interface{}{targetValue}, value)
		if err := m.Delete(ctx, "key"); err != nil {
			t.Fatal(err)
		}
		if value := it.MustSlice(m.Get(ctx, "key")); len(value) != 0 {
			t.Fatalf("target nil != %v", value)
		}
	})
}

func TestMultiMap_Clear(t *testing.T) {
	it.MultiMapTester(t, func(t *testing.T, m *hz.MultiMap) {
		ctx := context.Background()
		targetValue := "value"
		assert.True(t, it.MustBool(m.Put(ctx, "key", targetValue)), "multi-map put failed")
		value := it.MustValue(m.Get(ctx, "key"))
		it.AssertEquals(t, []interface{}{targetValue}, value) // check value
		if err := m.Clear(ctx); err != nil {
			t.Fatal(err)
		}
		if value := it.MustSlice(m.Get(ctx, "key")); len(value) != 0 {
			t.Fatalf("target nil != %v", value)
		}
	})
}

func TestMultiMap_Remove(t *testing.T) {
	it.MultiMapTester(t, func(t *testing.T, m *hz.MultiMap) {
		ctx := context.Background()
		targetValue := "value"
		assert.True(t, it.MustBool(m.Put(ctx, "key", targetValue)), "multi-map put failed")
		value := it.MustValue(m.Get(ctx, "key"))
		it.AssertEquals(t, []interface{}{targetValue}, value) // check value
		values, err := m.Remove(ctx, "key")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, []interface{}{targetValue}, values)
		if values = it.MustSlice(m.Get(ctx, "key")); len(values) != 0 {
			t.Fatalf("target nil != %v", values)
		}
	})
}

func TestMultiMap_GetKeySet(t *testing.T) {
	it.MultiMapTester(t, func(t *testing.T, m *hz.MultiMap) {
		ctx := context.Background()
		targetKeySet := []interface{}{"k1", "k2", "k3"}
		targets := []types.Entry{
			types.NewEntry("k1", "v1"),
			types.NewEntry("k2", "v2"),
			types.NewEntry("k3", "v3"),
		}
		for _, target := range targets {
			assert.True(t, it.MustBool(m.Put(ctx, target.Key, target.Value)), "multi-map put failed")
		}
		keySet, err := m.GetKeySet(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.ElementsMatch(t, targetKeySet, keySet)
	})
}

func TestMultiMap_GetValues(t *testing.T) {
	it.MultiMapTester(t, func(t *testing.T, m *hz.MultiMap) {
		ctx := context.Background()
		targetValues := []interface{}{"v1", "v2", "v3"}
		targets := []types.Entry{
			types.NewEntry("k1", "v1"),
			types.NewEntry("k2", "v2"),
			types.NewEntry("k3", "v3"),
		}
		for _, target := range targets {
			assert.True(t, it.MustBool(m.Put(ctx, target.Key, target.Value)), "multi-map put failed")
		}
		for _, target := range targets {
			it.AssertEquals(t, []interface{}{target.Value}, it.MustValue(m.Get(ctx, target.Key)))
		}
		values, err := m.GetValues(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.ElementsMatch(t, targetValues, values)
	})
}

func TestMultiMap_GetEntrySet(t *testing.T) {
	it.MultiMapTester(t, func(t *testing.T, m *hz.MultiMap) {
		ctx := context.Background()
		targets := []types.Entry{
			types.NewEntry("k1", "v1"),
			types.NewEntry("k2", "v2"),
			types.NewEntry("k3", "v3"),
		}
		for _, target := range targets {
			assert.True(t, it.MustBool(m.Put(ctx, target.Key, target.Value)), "multi-map put failed")
		}
		entries, err := m.GetEntrySet(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.ElementsMatch(t, targets, entries)
	})
}

func TestMultiMap_Lock(t *testing.T) {
	it.MultiMapTester(t, func(t *testing.T, cm *hz.MultiMap) {
		ctx := context.Background()
		const goroutineCount = 100
		const key = "counter"
		wg := &sync.WaitGroup{}
		wg.Add(goroutineCount)
		for i := 0; i < goroutineCount; i++ {
			go func() {
				defer wg.Done()
				intValue := int64(0)
				lockCtx := cm.NewLockContext(ctx)
				if err := cm.Lock(lockCtx, key); err != nil {
					panic(err)
				}
				defer cm.Unlock(lockCtx, key)
				v, err := cm.Get(lockCtx, key)
				if err != nil {
					panic(err)
				}
				if len(v) != 0 {
					intValue = v[0].(int64)
				}
				intValue++
				if _, err := cm.Remove(lockCtx, key); err != nil {
					panic(err)
				}
				successful, err := cm.Put(lockCtx, key, intValue)
				if err != nil {
					panic(err)
				}
				if !successful {
					panic("operation multi-map put failed")
				}
			}()
		}
		wg.Wait()
		it.AssertEquals(t, []interface{}{int64(goroutineCount)}, it.MustValue(cm.Get(ctx, key)))
	})
}

func TestMultiMap_ForceUnlockLock(t *testing.T) {
	it.MultiMapTester(t, func(t *testing.T, cm *hz.MultiMap) {
		ctx := context.Background()
		lockCtx := cm.NewLockContext(ctx)
		if err := cm.Lock(lockCtx, "k1"); err != nil {
			t.Fatal(err)
		}
		locked, err := cm.IsLocked(lockCtx, "k1")
		if err != nil {
			t.Fatal(err)
		}
		it.AssertEquals(t, true, locked)
		if err := cm.ForceUnlock(lockCtx, "k1"); err != nil {
			t.Fatal(err)
		}
		locked, err = cm.IsLocked(lockCtx, "k1")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, false, locked)
	})
}

func TestMultiMap_Size(t *testing.T) {
	it.MultiMapTester(t, func(t *testing.T, m *hz.MultiMap) {
		ctx := context.Background()
		targetSize := 0
		value, err := m.Size(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if targetSize != value {
			t.Fatalf("target: %d != %d", targetSize, value)
		}
		it.MustValue(m.Put(ctx, "k1", "v1"))
		it.MustValue(m.Put(ctx, "k2", "v2"))
		it.MustValue(m.Put(ctx, "k3", "v3"))
		targetSize = 3
		it.Eventually(t, func() bool {
			value, err := m.Size(ctx)
			if err != nil {
				t.Fatal(err)
			}
			return targetSize == value
		})
	})
}

func TestMultiMap_EntryNotifiedEvent(t *testing.T) {
	it.MultiMapTester(t, func(t *testing.T, m *hz.MultiMap) {
		ctx := context.Background()
		const totalCallCount = int32(100)
		callCount := int32(0)
		handler := func(event *hz.EntryNotified) {
			if event.EventType == hz.EntryAdded {
				atomic.AddInt32(&callCount, 1)
			}
		}
		listenerConfig := hz.MultiMapEntryListenerConfig{
			IncludeValue: true,
		}
		listenerConfig.NotifyEntryAdded(true)
		subscriptionID, err := m.AddEntryListener(ctx, listenerConfig, handler)
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < int(totalCallCount); i++ {
			key := fmt.Sprintf("key-%d", i)
			value := fmt.Sprintf("value-%d", i)
			it.MustValue(m.Put(ctx, key, value))
		}
		it.Eventually(t, func() bool {
			return atomic.LoadInt32(&callCount) == totalCallCount
		})
		atomic.StoreInt32(&callCount, 0)
		if err := m.RemoveEntryListener(ctx, subscriptionID); err != nil {
			t.Fatal(err)
		}
		for i := 0; i < int(totalCallCount); i++ {
			key := fmt.Sprintf("key-%d", i)
			value := fmt.Sprintf("value-%d", i)
			it.MustValue(m.Put(ctx, key, value))
		}
		if !assert.Equal(t, int32(0), atomic.LoadInt32(&callCount)) {
			t.FailNow()
		}
	})
}

func TestMultiMap_EntryNotifiedEventToKey(t *testing.T) {
	it.MultiMapTester(t, func(t *testing.T, m *hz.MultiMap) {
		ctx := context.Background()
		callCount := int32(0)
		handler := func(event *hz.EntryNotified) {
			atomic.AddInt32(&callCount, 1)
		}
		listenerConfig := hz.MultiMapEntryListenerConfig{
			IncludeValue: true,
			Key:          "k1",
		}
		listenerConfig.NotifyEntryAdded(true)
		if _, err := m.AddEntryListener(ctx, listenerConfig, handler); err != nil {
			t.Fatal(err)
		}
		it.MustValue(m.Put(ctx, "k1", "v1"))
		it.Eventually(t, func() bool {
			return atomic.LoadInt32(&callCount) == int32(1)
		})
	})
}

func TestMultiMap_Destroy(t *testing.T) {
	it.MultiMapTester(t, func(t *testing.T, m *hz.MultiMap) {
		ctx := context.Background()
		if err := m.Destroy(ctx); err != nil {
			t.Fatal(err)
		}
	})
}

func TestMultiMap_TryLock(t *testing.T) {
	it.MultiMapTester(t, func(t *testing.T, m *hz.MultiMap) {
		ctx := context.Background()
		wg := &sync.WaitGroup{}
		wg.Add(1)
		const key = "foo"
		go func() {
			ctx := m.NewLockContext(ctx)
			it.Must(m.Lock(ctx, key))
			wg.Done()
		}()
		wg.Wait()
		mainCtx := m.NewLockContext(ctx)
		b, err := m.TryLock(mainCtx, key)
		if err != nil {
			t.Fatal(err)
		}
		assert.False(t, b)
	})
}

func TestMultiMap_LockWithLease(t *testing.T) {
	it.MultiMapTester(t, func(t *testing.T, m *hz.MultiMap) {
		ctx := context.Background()
		wg := &sync.WaitGroup{}
		wg.Add(1)
		const key = "foo"
		go func() {
			ctx := m.NewLockContext(ctx)
			it.Must(m.LockWithLease(ctx, key, 100*time.Millisecond))
			wg.Done()
		}()
		wg.Wait()
		mainCtx := m.NewLockContext(ctx)
		err := m.Lock(mainCtx, key)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestMultiMap_TryLockWithLease(t *testing.T) {
	it.MultiMapTester(t, func(t *testing.T, m *hz.MultiMap) {
		ctx := context.Background()
		wg := &sync.WaitGroup{}
		wg.Add(1)
		const key = "foo"
		go func() {
			ctx := m.NewLockContext(ctx)
			if b := it.MustValue(m.TryLockWithLease(ctx, key, 100*time.Millisecond)); b != true {
				panic("unexpected value")
			}
			wg.Done()
		}()
		wg.Wait()
		mainCtx := m.NewLockContext(ctx)
		err := m.Lock(mainCtx, key)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestMultiMap_TryLockWithTimeout(t *testing.T) {
	it.MultiMapTester(t, func(t *testing.T, m *hz.MultiMap) {
		ctx := context.Background()
		wg := &sync.WaitGroup{}
		wg.Add(1)
		const key = "foo"
		go func() {
			ctx := m.NewLockContext(ctx)
			it.Must(m.Lock(ctx, key))
			wg.Done()
		}()
		wg.Wait()
		mainCtx := m.NewLockContext(ctx)
		b, err := m.TryLockWithTimeout(mainCtx, key, 100*time.Millisecond)
		if err != nil {
			t.Fatal(err)
		}
		assert.False(t, b)
	})
}

func TestMultiMap_TryLockWithLeaseAndTimeout(t *testing.T) {
	it.MultiMapTester(t, func(t *testing.T, m *hz.MultiMap) {
		ctx := context.Background()
		wg := &sync.WaitGroup{}
		wg.Add(1)
		const key = "foo"
		go func() {
			ctx := m.NewLockContext(ctx)
			it.Must(m.Lock(ctx, key))
			wg.Done()
		}()
		wg.Wait()
		mainCtx := m.NewLockContext(ctx)
		b, err := m.TryLockWithLeaseAndTimeout(mainCtx, key, 100*time.Millisecond, 200*time.Millisecond)
		if err != nil {
			t.Fatal(err)
		}
		assert.False(t, b)
	})
}
