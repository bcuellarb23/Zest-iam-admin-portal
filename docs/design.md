Developer: Bryan Consuegra Cuellar

# DESIGN.md - Project Architecture and Rules

## 1. Overview
Goal 
The goal of this update is to build a concurrency-safe in-memory session and token store in Go. This Sprint is to rewrite the IAM session state management or authentication flow in go. 

## 2. Package Rules
/cmd 
will house the main application entry points

/internal 
will have my private application code such as auth logic, and state store

/pkg
This folder will have code that is ok to use by external applications, other projects will import these libraries expecting them to work.

## 3. Concurrency model
To manage concurrent access to active user sessions and roles in memory, this service will use shared memory with mutex
- ** Rationale:** Incoming HTTP requests will frequently read session states and ocasionally write/update tokens during login/logout/ A 'sync.RWMutex' permits concurrent read operation for optimal performance while safely locking the data store during write operations to prevent data races.

- ** Race Condition Prevention:** All state mutations will be guarded by write locks ('Lock()'), and lookups will use read locks ('RLock()'). Race detection will be enforced during testing via 'go test -race'.


```

